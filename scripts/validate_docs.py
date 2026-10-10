#!/usr/bin/env python3
"""Validate CI See documentation navigation and stable contract identifiers.

Stdlib only. Historical migration hashes are evidence, not a lock on future edits.
This cannot certify semantic equivalence, implementation behavior or hosted parity.
"""
import argparse
from collections import Counter
import json
from pathlib import Path
import re
import sys
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / ".agents/skills/context-engineering/scripts"))
from inspect_docs import inspect, markdown_links, visible_markdown  # noqa: E402


def slugify(title):
    return re.sub(r"[^\w\s-]", "", title.lower()).strip().replace(" ", "-")


def anchors(content):
    titles = re.findall(r"(?m)^ {0,3}#{1,6}[ \t]+(.+?)\s*$",
                        visible_markdown(content))
    counts = Counter()
    found = set()
    for title in titles:
        slug = slugify(title.rstrip(" #"))
        suffix = counts[slug]
        counts[slug] += 1
        found.add(slug if suffix == 0 else f"{slug}-{suffix}")
    return found


def contract_ids(root):
    failures = []
    prd_parts = list((root / "docs/prd").glob("*.md"))
    fr = []
    for path in prd_parts:
        fr.extend(int(x) for x in re.findall(
            r"(?m)^### FR-(\d+)\b", path.read_text(encoding="utf-8")))
    if sorted(fr) != list(range(1, 14)):
        failures.append(f"FR definitions must uniquely cover 1..13; got {sorted(fr)}")

    tests = [
        ("local-ci-execution", 12),
        ("git-gating", 20),
        ("developer-approval-ux", 23),
    ]
    for spec, maximum in tests:
        src_dir = root / "docs/specs" / spec
        all_text = "\n".join(p.read_text(encoding="utf-8")
                             for p in sorted(src_dir.glob("*.md")))
        matches = re.findall(r"(?m)^## \d+\. Required Acceptance Tests\s*$",
                             all_text)
        if len(matches) != 1:
            failures.append(f"{spec}: expected one acceptance section, got {len(matches)}")
            continue
        section = re.split(r"(?m)^## \d+\. Required Acceptance Tests\s*$",
                           all_text, maxsplit=1)[1]
        section = re.split(r"(?m)^## ", section, maxsplit=1)[0]
        numbers = [int(n) for n in re.findall(r"(?m)^(\d+)\. ", section)]
        if numbers != list(range(1, maximum + 1)):
            failures.append(f"{spec}: acceptance numbers {numbers} != 1..{maximum}")

    plan = root / "docs/implementation-plan"
    slices = (plan / "vertical-slices.md").read_text(encoding="utf-8")
    gates = (plan / "authority-and-gates.md").read_text(encoding="utf-8")
    s_ids = [int(n) for n in re.findall(r"(?m)^\|\s*\*\*S(\d+)\s+—", slices)]
    d_ids = [int(n) for n in re.findall(r"(?m)^\|\s*D(\d+)\s+—", gates)]
    if sorted(s_ids) != list(range(11)):
        failures.append(f"S slices are not unique 0..10: {s_ids}")
    if sorted(d_ids) != list(range(1, 5)):
        failures.append(f"D gates are not unique 1..4: {d_ids}")
    return failures


def verify(root):
    audit = inspect(root, ["docs"], ["docs/index.md"], [], 150)
    failures = [f"broken path: {x}" for x in audit["broken_links"]]
    failures += [f"unindexed: {p}" for p in audit["unindexed_documents"]]
    failures += [f"missing entry: {p}" for p in audit["missing_entry_points"]]
    failures += [f"skipped: {p}" for p in audit["skipped_symlink_files"]]
    failures += [f"missing docs root: {p}" for p in audit["missing_doc_roots"]]
    file_cache = {}

    def read(path):
        if path not in file_cache:
            file_cache[path] = path.read_text(encoding="utf-8")
        return file_cache[path]

    for row in audit["documents"]:
        path = root / row["path"]
        for line, url in markdown_links(read(path)):
            parsed = urlsplit(url)
            if parsed.scheme or parsed.netloc or not parsed.fragment:
                continue
            target = (path.parent / unquote(parsed.path)).resolve() if parsed.path else path
            if not target.is_file() or not target.is_relative_to(root):
                continue  # missing-path finding already reported by bundled scanner
            fragment = unquote(parsed.fragment)
            if fragment not in anchors(read(target)):
                failures.append(f"broken fragment: {row['path']}:{line} -> {url}")

    manifest = json.loads(read(root / "docs/migration/manifest.json"))
    for source in manifest["sources"]:
        stub = root / source["source"]
        content = read(stub)
        original_anchors = anchors(content)
        for anchor in source["anchors"]:
            fragment = anchor["legacy_fragment"]
            if fragment not in original_anchors:
                failures.append(f"lost legacy anchor: {source['source']}#{fragment}")
            dest = root / anchor["destination"]
            # Ensure the legacy heading maps to its exact destination fragment.
            expected_rel = __import__("os").path.relpath(dest, stub.parent).replace("\\", "/")
            expected_link = f"({expected_rel}#{anchor['destination_fragment']})"
            if expected_link not in content:
                failures.append(f"legacy heading route missing: {source['source']}#{fragment}")
            if anchor["destination_fragment"] not in anchors(read(dest)):
                failures.append(f"legacy destination anchor missing: {dest}#{fragment}")
        for child in source["children"]:
            if not (root / child["path"]).is_file():
                failures.append(f"missing migrated section: {child['path']}")
    failures += contract_ids(root)
    return audit, failures


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true",
                        help="nonzero exit on navigation or contract-ID failures")
    args = parser.parse_args()
    audit, failures = verify(ROOT)
    report = {
        "status": "FAIL" if failures else "PASS",
        "files": len(audit["documents"]),
        "broken_paths": len(audit["broken_links"]),
        "unindexed": len(audit["unindexed_documents"]),
        "over_150_review": [d["path"] for d in audit["documents"]
                            if d["review_required"]],
        "sizes": {d["path"]: d["lines"] for d in audit["documents"]},
        "failures": failures,
        "limits": "Does not prove semantic equivalence, runtime correctness or external link availability",
    }
    print(json.dumps(report, indent=2))
    return int(args.check and bool(failures))


if __name__ == "__main__":
    raise SystemExit(main())
