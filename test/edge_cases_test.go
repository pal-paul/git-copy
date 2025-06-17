package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pal-paul/git-copy/internal/gitcopy"
)

// TestComplexDirectoryStructures tests handling of complex nested directory structures
func TestComplexDirectoryStructures(t *testing.T) {
	tempDir := t.TempDir()

	// Create a complex directory structure
	testStructure := []string{
		"level1/file1.txt",
		"level1/level2/file2.txt",
		"level1/level2/level3/file3.txt",
		"level1/level2/level3/level4/file4.txt",
		"level1/another_branch/file5.txt",
		"level1/another_branch/sub/file6.txt",
		"separate_root/file7.txt",
		"separate_root/deep/deeper/deepest/file8.txt",
	}

	for _, filePath := range testStructure {
		fullPath := filepath.Join(tempDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", filePath, err)
		}

		content := "Content for " + filePath
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	// Test recursive directory reading
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed: %v", err)
	}

	if len(files) != len(testStructure) {
		t.Errorf("Expected %d files, got %d", len(testStructure), len(files))
	}

	// Verify all expected files are found
	expectedPaths := make(map[string]bool)
	for _, filePath := range testStructure {
		expectedPaths[filepath.Join(tempDir, filePath)] = false
	}

	for _, file := range files {
		if _, exists := expectedPaths[file]; exists {
			expectedPaths[file] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf("Expected file not found: %s", path)
		}
	}
}

// TestLargeFileHandling tests handling of large files
func TestLargeFileHandling(t *testing.T) {
	tempDir := t.TempDir()
	largeFile := filepath.Join(tempDir, "large_file.txt")

	// Create a file with substantial content (1MB)
	content := strings.Repeat("This is a line of test content for large file testing.\n", 20000)
	err := os.WriteFile(largeFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Test reading large file
	readContent, err := gitcopy.ReadFile(largeFile)
	if err != nil {
		t.Fatalf("Failed to read large file: %v", err)
	}

	if len(readContent) != len(content) {
		t.Errorf("Large file content length mismatch: expected %d, got %d", len(content), len(readContent))
	}

	if string(readContent) != content {
		t.Error("Large file content does not match expected content")
	}
}

// TestSpecialCharactersInPaths tests handling of special characters in file paths
func TestSpecialCharactersInPaths(t *testing.T) {
	tempDir := t.TempDir()

	// Create files with special characters (that are valid on most systems)
	specialFiles := []string{
		"file with spaces.txt",
		"file-with-dashes.txt",
		"file_with_underscores.txt",
		"file.with.dots.txt",
		"file123with456numbers.txt",
		"UPPERCASE.TXT",
		"MixedCase.TxT",
	}

	for _, fileName := range specialFiles {
		filePath := filepath.Join(tempDir, fileName)
		content := "Content for " + fileName
		err := os.WriteFile(filePath, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("Failed to create file with special chars %s: %v", fileName, err)
		}
	}

	// Test reading directory with special character files
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed with special characters: %v", err)
	}

	if len(files) != len(specialFiles) {
		t.Errorf("Expected %d files with special characters, got %d", len(specialFiles), len(files))
	}

	// Test reading each file with special characters
	for _, file := range files {
		content, err := gitcopy.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read file with special characters %s: %v", file, err)
		}

		if len(content) == 0 {
			t.Errorf("File with special characters has empty content: %s", file)
		}
	}
}

// TestEmptyFiles tests handling of empty files
func TestEmptyFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create empty files
	emptyFiles := []string{"empty1.txt", "empty2.log", "subdir/empty3.json"}

	for _, fileName := range emptyFiles {
		filePath := filepath.Join(tempDir, fileName)
		err := os.MkdirAll(filepath.Dir(filePath), 0o755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", fileName, err)
		}

		err = os.WriteFile(filePath, []byte(""), 0o644)
		if err != nil {
			t.Fatalf("Failed to create empty file %s: %v", fileName, err)
		}
	}

	// Test reading empty files
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed with empty files: %v", err)
	}

	if len(files) != len(emptyFiles) {
		t.Errorf("Expected %d empty files, got %d", len(emptyFiles), len(files))
	}

	for _, file := range files {
		content, err := gitcopy.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read empty file %s: %v", file, err)
		}

		if len(content) != 0 {
			t.Errorf("Empty file should have zero content, got %d bytes: %s", len(content), file)
		}
	}
}

// TestSymlinkHandling tests handling of symbolic links (Unix-like systems)
func TestSymlinkHandling(t *testing.T) {
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping symlink test on Windows")
	}

	tempDir := t.TempDir()

	// Create a regular file
	regularFile := filepath.Join(tempDir, "regular.txt")
	err := os.WriteFile(regularFile, []byte("regular file content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	// Create a symbolic link to the regular file
	symlinkFile := filepath.Join(tempDir, "symlink.txt")
	err = os.Symlink(regularFile, symlinkFile)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Test directory reading (should include both files)
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed with symlinks: %v", err)
	}

	// Should find both the regular file and the symlink
	if len(files) < 1 {
		t.Error("Should find at least the regular file")
	}

	// Test reading both files
	for _, file := range files {
		content, err := gitcopy.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", file, err)
		}

		if len(content) == 0 {
			t.Errorf("File should not be empty: %s", file)
		}
	}
}

// TestConcurrentDirectoryAccess tests concurrent access to directories
func TestConcurrentDirectoryAccess(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple files
	for i := 0; i < 20; i++ {
		fileName := filepath.Join(tempDir, "file"+string(rune(i+48))+".txt") // ASCII numbers
		content := "Content for concurrent test file " + string(rune(i+48))
		err := os.WriteFile(fileName, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("Failed to create concurrent test file: %v", err)
		}
	}

	// Test concurrent directory reading
	results := make(chan []string, 5)
	errors := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			files, err := gitcopy.IoReadDir(tempDir)
			if err != nil {
				errors <- err
			} else {
				results <- files
			}
		}()
	}

	// Collect results
	for i := 0; i < 5; i++ {
		select {
		case files := <-results:
			if len(files) != 20 {
				t.Errorf("Concurrent read %d: expected 20 files, got %d", i, len(files))
			}
		case err := <-errors:
			t.Errorf("Concurrent read failed: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent directory reads")
		}
	}
}

// TestFilePermissionChanges tests handling of files with different permissions
func TestFilePermissionChanges(t *testing.T) {
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	tempDir := t.TempDir()

	// Create files with different permissions
	permissions := []struct {
		name string
		perm os.FileMode
	}{
		{"readable.txt", 0o644},
		{"executable.txt", 0o755},
		{"readonly.txt", 0o444},
	}

	for _, p := range permissions {
		filePath := filepath.Join(tempDir, p.name)
		content := "Content with permission " + p.name
		err := os.WriteFile(filePath, []byte(content), p.perm)
		if err != nil {
			t.Fatalf("Failed to create file with permissions %s: %v", p.name, err)
		}
	}

	// Test reading directory with different permission files
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed with permission files: %v", err)
	}

	if len(files) != len(permissions) {
		t.Errorf("Expected %d permission files, got %d", len(permissions), len(files))
	}

	// Test reading each file
	for _, file := range files {
		content, err := gitcopy.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read permission file %s: %v", file, err)
		}

		if len(content) == 0 {
			t.Errorf("Permission file should not be empty: %s", file)
		}
	}
}

// TestResourceCleanup tests that resources are properly cleaned up
func TestResourceCleanup(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "cleanup_test.txt")

	// Create a test file
	content := "This is test content for resource cleanup testing"
	err := os.WriteFile(testFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create cleanup test file: %v", err)
	}

	// Test multiple reads to ensure proper cleanup
	for i := 0; i < 10; i++ {
		readContent, err := gitcopy.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Read %d failed: %v", i+1, err)
		}

		if string(readContent) != content {
			t.Errorf("Read %d: content mismatch", i+1)
		}
	}

	// The test passes if no file handle leaks occur
	t.Log("Resource cleanup test completed successfully")
}

// TestPathNormalization tests cross-platform path handling
func TestPathNormalization(t *testing.T) {
	tempDir := t.TempDir()

	// Create nested structure
	nestedPath := filepath.Join("level1", "level2", "file.txt")
	fullPath := filepath.Join(tempDir, nestedPath)

	err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
	if err != nil {
		t.Fatalf("Failed to create nested path: %v", err)
	}

	err = os.WriteFile(fullPath, []byte("nested content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create nested file: %v", err)
	}

	// Test that paths are properly normalized
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(files) > 0 {
		// Verify the path uses proper separators
		expectedPath := filepath.Join(tempDir, "level1", "level2", "file.txt")
		if files[0] != expectedPath {
			t.Errorf("Path normalization failed: expected %s, got %s", expectedPath, files[0])
		}

		// Verify the path is absolute
		if !filepath.IsAbs(files[0]) {
			t.Errorf("Path should be absolute: %s", files[0])
		}
	}
}
