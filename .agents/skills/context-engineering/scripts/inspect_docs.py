#!/usr/bin/env python3
"""Inventory selected Markdown and check common local file links. Stdlib only.

Emit JSON; never rewrite documents. Fragment anchors, HTML, generated routing,
semantic preservation, and implementation correctness require separate checks.
"""

import argparse
import fnmatch
import hashlib
import json
import os
from pathlib import Path
import re
from urllib.parse import unquote, urlsplit


PRUNED = {".git", "node_modules", ".next", ".venv", "__pycache__"}


def visible_markdown(content, *, mask_inline=False):
    """Mask common code examples and comments without moving source positions.

    Process constructs in source order: a comment marker in code is literal,
    and a fence inside a comment cannot open a code block. This is a lightweight
    scanner for top-level Markdown, not a complete CommonMark block parser.
    """
    result = list(content)

    def mask(start, end):
        for position in range(start, end):
            if content[position] not in "\r\n":
                result[position] = " "

    position, fence, indented, paragraph = 0, None, False, False
    while position < len(content):
        if position == 0 or content[position - 1] == "\n":
            end = content.find("\n", position)
            end = len(content) if end == -1 else end + 1
            line = content[position:end]
            if fence:
                if re.fullmatch(r" {0,3}" + re.escape(fence[0]) +
                                r"{" + str(len(fence)) + r",}\s*", line):
                    fence = None
                mask(position, end)
                position = end
                continue
            is_indented = line.expandtabs(4).startswith("    ")
            if (indented and (is_indented or not line.strip()) or
                    is_indented and not paragraph):
                indented, paragraph = True, False
                mask(position, end)
                position = end
                continue
            indented = False
            opening = re.match(r" {0,3}(`{3,}|~{3,})(.*)", line)
            if opening and not (opening[1][0] == "`" and "`" in opening[2]):
                fence, paragraph = opening[1], False
                mask(position, end)
                position = end
                continue
            paragraph = bool(line.strip()) and not re.match(
                r" {0,3}(?:#{1,6}(?:[ \t]|$)|(?:[*_-][ \t]*){3,}$)", line)

        if content[position] == "\\":
            # Escaped backticks/comment openers are ordinary text.
            position += 2 if position + 1 < len(content) else 1
            continue
        if content.startswith("<!--", position):
            close = content.find("-->", position + 4)
            end = len(content) if close == -1 else close + 3
            line_start = content.rfind("\n", 0, position) + 1
            line_end = content.find("\n", end)
            line_end = len(content) if line_end == -1 else line_end
            if (not content[line_start:position].strip() and
                    not content[end:line_end].strip()):
                paragraph = False
            mask(position, end)
            position = end
            continue
        if content[position] == "`":
            run = re.match(r"`+", content[position:])[0]
            start = position + len(run)
            # Inline code may cross lines, but cannot cross a paragraph/block.
            boundary = re.search(r"\n[ \t]*\n|\n {0,3}(?:`{3,}|~{3,})", content[start:])
            limit = start + boundary.start() if boundary else len(content)
            close = re.search(r"(?<!`)" + re.escape(run) + r"(?!`)",
                              content[start:limit])
            if close:
                end = start + close.end()
                if mask_inline:
                    mask(position, end)
                position = end
            else:
                position = start
            continue
        position += 1
    return "".join(result)


def destination(text, start):
    """Read a common inline/reference destination, allowing nested parentheses."""
    while start < len(text) and text[start].isspace():
        start += 1
    if start < len(text) and text[start] == "<":
        end = text.find(">", start + 1)
        return text[start + 1:end] if end != -1 else None
    value, depth, i = [], 0, start
    while i < len(text):
        char = text[i]
        if char == "\\" and i + 1 < len(text):
            value.append(text[i + 1])
            i += 2
            continue
        if char == "(":
            depth += 1
        elif char == ")":
            if depth == 0:
                break
            depth -= 1
        elif char.isspace() and depth == 0:
            break
        value.append(char)
        i += 1
    return "".join(value) if depth == 0 else None


def normalize_label(label):
    return " ".join(label.split()).casefold()


def markdown_links(text):
    text = visible_markdown(text, mask_inline=True)
    definitions = {}
    for match in re.finditer(r"(?m)^ {0,3}\[([^\]\n]+)\]:[ \t]*", text):
        target = destination(text, match.end())
        if target is not None:
            definitions[normalize_label(match.group(1))] = target
    links = set()
    for match in re.finditer(r"!?\[[^\]\n]*\]\(", text):
        target = destination(text, match.end())
        if target is not None:
            links.add((text.count("\n", 0, match.start()) + 1, target))
    for match in re.finditer(r"!?\[([^\]\n]+)\](?:\[([^\]\n]*)\])?", text):
        if text[match.end():match.end() + 1] in {"(", ":"}:
            continue
        label = normalize_label(match.group(2) or match.group(1))
        if label in definitions:
            links.add((text.count("\n", 0, match.start()) + 1, definitions[label]))
    return sorted(links)


def inside(path, root):
    try:
        path.relative_to(root)
        return True
    except ValueError:
        return False


def raise_walk_error(error):
    raise error


def inspect(root, doc_roots, entries, exclusions, threshold):
    selected = {root / name for name in ("README.md", "AGENTS.md")
                if (root / name).is_file()}
    missing_roots, skipped, skipped_dirs = [], [], []
    for name in doc_roots:
        directory = root / name
        if not directory.is_dir():
            missing_roots.append(name)
            continue
        if directory.is_symlink():
            skipped_dirs.append(name)
            continue
        for current, dirs, files in os.walk(directory, followlinks=False,
                                           onerror=raise_walk_error):
            skipped_dirs.extend((Path(current) / d).relative_to(root).as_posix()
                                for d in dirs if (Path(current) / d).is_symlink())
            dirs[:] = sorted(d for d in dirs if d not in PRUNED and
                             not (Path(current) / d).is_symlink())
            selected.update(Path(current) / f for f in files if f.endswith(".md"))
    selected = {p for p in selected
                if not any(fnmatch.fnmatch(p.relative_to(root).as_posix(), pattern)
                           for pattern in exclusions)}
    documents, edges, broken = [], {}, []
    for path in sorted(selected):
        rel = path.relative_to(root).as_posix()
        if path.is_symlink() or not inside(path.resolve(), root):
            skipped.append(rel)
            continue
        raw = path.read_bytes()
        content = raw.decode("utf-8-sig")
        visible = visible_markdown(content)
        headings = [{"line": i, "title": m.group(1).rstrip(" #")}
                    for i, line in enumerate(visible.splitlines(), 1)
                    if (m := re.match(r" {0,3}#{1,6}[ \t]+(.+)", line))]
        intro = re.split(r"(?m)^ {0,3}#{2,6}[ \t]+", visible, maxsplit=1)[0]
        metadata = {m.group(1).lower(): m.group(2).strip()
                    for m in re.finditer(r"(?mi)^(Status|Authority):[ \t]*(.*)$", intro)}
        edges[rel] = set()
        for line, target in markdown_links(content):
            parsed = urlsplit(target)
            if parsed.scheme or parsed.netloc:
                continue
            link_path = unquote(parsed.path)
            resolved = ((root / link_path.lstrip("/")) if link_path.startswith("/")
                        else path.parent / link_path).resolve() if link_path else path
            if not inside(resolved, root):
                reason = "outside_repository"
            elif not resolved.exists():
                reason = "missing_target"
            else:
                edges[rel].add(resolved.relative_to(root).as_posix())
                continue
            broken.append({"source": rel, "line": line, "target": target,
                           "reason": reason})
        documents.append({"path": rel, "lines": len(content.splitlines()),
                          "bytes": len(raw), "sha256": hashlib.sha256(raw).hexdigest(),
                          "review_required": len(content.splitlines()) > threshold,
                          "headings": headings, "metadata": metadata,
                          "local_targets": sorted(edges[rel])})
    found = {d["path"] for d in documents}
    explicit = entries is not None
    entries = entries if explicit else [str(Path(n) / "index.md") for n in doc_roots
                                        if (root / n / "index.md").is_file()]
    missing_entries = sorted(e for e in entries if e not in found)
    reached, pending = set(), list(entries)
    while pending:
        current = pending.pop()
        if current not in reached:
            reached.add(current)
            pending.extend(edges.get(current, set()) - reached)
    unindexed = sorted(p for p in found - reached
                       if p not in {"AGENTS.md", "README.md"}) if entries else []
    return {"schema_version": 1, "root": str(root), "review_threshold": threshold,
            "doc_roots": doc_roots, "exclusions": exclusions, "documents": documents,
            "broken_links": broken, "entry_points": entries,
            "missing_entry_points": missing_entries, "missing_doc_roots": missing_roots,
            "index_coverage_checked": bool(entries), "unindexed_documents": unindexed,
            "skipped_symlink_files": skipped,
            "skipped_symlink_directories": sorted(d for d in skipped_dirs if not any(
                fnmatch.fnmatch(d, p) or fnmatch.fnmatch(d + "/", p) for p in exclusions)),
            "limitations": ["Common Markdown file links only; HTML and generated routing excluded.",
                            "Code/comment masking covers top-level constructs; nested list/blockquote parsing is not supported.",
                            "Fragment anchors and setext headings are not checked.",
                            "Symlink directories and standard dependency directories are not scanned.",
                            "Counts and hashes do not prove semantic equivalence or implementation freshness."]}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("root", type=Path)
    parser.add_argument("--docs", action="append", help="Relative documentation root; repeatable")
    parser.add_argument("--entry", action="append", help="Relative index entry point; repeatable")
    parser.add_argument("--exclude", action="append", default=[], help="Repo-relative glob; repeatable")
    parser.add_argument("--threshold", type=int, default=150)
    parser.add_argument("--check", action="store_true", help="Fail on broken links or index coverage")
    args = parser.parse_args()
    root = args.root.resolve()
    if not root.is_dir() or args.threshold < 1:
        parser.error("root must be a directory and threshold must be positive")
    doc_roots = args.docs or ["docs"]
    for name in doc_roots + (args.entry or []):
        if Path(name).is_absolute() or not inside((root / name).resolve(), root):
            parser.error("documentation roots and entry points must be inside the repository")
    result = inspect(root, doc_roots, args.entry, args.exclude, args.threshold)
    print(json.dumps(result, indent=2, ensure_ascii=False))
    failures = (result["broken_links"] or result["missing_entry_points"] or
                result["missing_doc_roots"] or result["unindexed_documents"] or
                result["skipped_symlink_files"] or result["skipped_symlink_directories"] or
                not result["index_coverage_checked"])
    return 1 if args.check and failures else 0


if __name__ == "__main__":
    raise SystemExit(main())
