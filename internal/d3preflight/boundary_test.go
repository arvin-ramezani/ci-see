package d3preflight

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProductionPackageHasNoExecutionOrRuntimeDependencies(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	productionDirectories := []string{workingDirectory, filepath.Join(workingDirectory, "..", "..", "cmd", "d3-preflight")}
	allowedImports := map[string]struct{}{
		"bytes": {}, "crypto/ed25519": {}, "crypto/sha256": {}, "encoding/base64": {}, "encoding/hex": {}, "encoding/json": {},
		"flag": {}, "fmt": {}, "io": {}, "net/netip": {}, "os": {}, "strings": {}, "time": {}, "unicode/utf8": {},
		"github.com/arvin-ramezani/ci-see/internal/d3preflight": {},
	}
	forbiddenCalls := map[string]struct{}{
		"Command": {}, "CommandContext": {}, "StartProcess": {}, "ForkExec": {}, "Exec": {}, "Dial": {}, "DialContext": {},
		"Listen": {}, "ListenPacket": {}, "ListenAndServe": {}, "Serve": {}, "Create": {}, "CreateTemp": {}, "Mkdir": {},
		"MkdirAll": {}, "Remove": {}, "RemoveAll": {}, "Rename": {}, "WriteFile": {}, "Chmod": {}, "Chown": {}, "Symlink": {}, "Link": {},
	}

	for _, directory := range productionDirectories {
		entries, err := os.ReadDir(directory)
		if err != nil {
			t.Fatalf("read production package directory: %v", err)
		}
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				continue
			}
			if !strings.HasSuffix(name, ".go") {
				switch strings.ToLower(filepath.Ext(name)) {
				case ".c", ".cc", ".cpp", ".h", ".s", ".syso":
					t.Errorf("%s contains native source or object file %q outside the Go import boundary", directory, name)
				}
				continue
			}
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			filename := filepath.Join(directory, name)
			file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.AllErrors)
			if err != nil {
				t.Fatalf("parse %s: %v", filename, err)
			}
			for _, imported := range file.Imports {
				importPath := strings.Trim(imported.Path.Value, "\"")
				if _, allowed := allowedImports[importPath]; !allowed {
					t.Errorf("%s imports out-of-bound package %q", filename, importPath)
				}
			}
			ast.Inspect(file, func(node ast.Node) bool {
				selector, ok := node.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if _, forbidden := forbiddenCalls[selector.Sel.Name]; forbidden {
					t.Errorf("%s uses forbidden process/network API selector %q", filename, selector.Sel.Name)
				}
				return true
			})
		}
	}
}

