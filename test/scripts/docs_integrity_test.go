package scripts_test

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var inlineMarkdownLink = regexp.MustCompile(`!?\[[^]]*\]\(([^)]+)\)`)

func TestDocumentationLinksResolve(t *testing.T) {
	t.Parallel()

	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Abs(repo root) error = %v", err)
	}

	files := []string{
		"README.md",
		"README_zh.md",
		"CONTRIBUTING.md",
		"scripts/README.md",
	}
	err = filepath.WalkDir(filepath.Join(root, "docs"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(path), ".md") {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(docs) error = %v", err)
	}

	sort.Strings(files)
	for _, rel := range files {
		rel := rel
		t.Run(filepath.ToSlash(rel), func(t *testing.T) {
			content, readErr := os.ReadFile(filepath.Join(root, rel))
			if readErr != nil {
				t.Fatalf("ReadFile() error = %v", readErr)
			}

			for _, match := range inlineMarkdownLink.FindAllStringSubmatch(string(content), -1) {
				target := strings.Trim(strings.TrimSpace(match[1]), "<>")
				if target == "" || strings.HasPrefix(target, "#") || hasExternalScheme(target) {
					continue
				}
				if separator := strings.IndexAny(target, "#?"); separator >= 0 {
					target = target[:separator]
				}
				if target == "" {
					continue
				}

				resolved := filepath.Clean(filepath.Join(root, filepath.Dir(rel), filepath.FromSlash(target)))
				if _, statErr := os.Stat(resolved); statErr != nil {
					t.Errorf("local link %q does not resolve: %v", match[1], statErr)
				}
			}
		})
	}
}

func TestEveryDocumentationPageIsIndexed(t *testing.T) {
	t.Parallel()

	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Abs(repo root) error = %v", err)
	}

	indexPath := filepath.Join(root, "docs", "README.md")
	index, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", indexPath, err)
	}

	err = filepath.WalkDir(filepath.Join(root, "docs"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") || path == indexPath {
			return nil
		}

		rel, relErr := filepath.Rel(filepath.Dir(indexPath), path)
		if relErr != nil {
			return relErr
		}
		target := "./" + filepath.ToSlash(rel)
		if !strings.Contains(string(index), "]("+target+")") {
			t.Errorf("%s is not linked from docs/README.md as %s", filepath.ToSlash(rel), target)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(docs) error = %v", err)
	}
}

func hasExternalScheme(target string) bool {
	lower := strings.ToLower(target)
	for _, prefix := range []string{"http://", "https://", "mailto:", "tel:", "data:", "app://"} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}
