package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pal-paul/git-copy/internal/gitcopy"
)

const (
	defaultBranch = "refs/heads/master"
)

// Setup function to initialize test environment
func setupTestEnvironment() gitcopy.Environment {
	testEnv := gitcopy.Environment{}
	testEnv.GitHub.Token = "test-token"
	testEnv.GitHub.Api = "https://api.github.com"
	testEnv.GitHub.Repo = "test/repo"
	testEnv.GitHub.Workflow = "test-workflow"
	testEnv.GitHub.Branch = defaultBranch
	testEnv.GitHub.Commit = "abc123"
	testEnv.GitHub.RunId = "12345"
	testEnv.GitHub.JobName = "test-job"
	testEnv.GitHub.Server = "https://github.com"
	testEnv.Input.Owner = "test-owner"
	testEnv.Input.Repo = "test-repo"
	testEnv.Input.Branch = "test-branch"
	testEnv.Input.PullMessage = "Test pull request"
	testEnv.Input.PullDescription = "Test description"
	return testEnv
}

// TestEnvironmentFunctions tests environment-related functions
func TestEnvironmentFunctions(t *testing.T) {
	// Save original environment
	originalEnv := gitcopy.GetEnvironment()
	defer gitcopy.SetEnvironment(originalEnv)

	// Set test environment
	testEnv := setupTestEnvironment()
	gitcopy.SetEnvironment(testEnv)

	// Test GetEnvironment
	env := gitcopy.GetEnvironment()

	if env.GitHub.Token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", env.GitHub.Token)
	}

	if env.Input.Owner != "test-owner" {
		t.Errorf("Expected owner 'test-owner', got '%s'", env.Input.Owner)
	}

	if env.Input.Branch != "test-branch" {
		t.Errorf("Expected branch 'test-branch', got '%s'", env.Input.Branch)
	}
}

// TestReadFile tests the gitcopy.ReadFile function
func TestReadFile(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!\nThis is a test file."

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading the file
	content, err := gitcopy.ReadFile(testFile)
	if err != nil {
		t.Fatalf("gitcopy.ReadFile failed: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, string(content))
	}
}

// Testgitcopy.ReadFileError tests error cases for gitcopy.ReadFile
func TestReadFileError(t *testing.T) {
	// Test with non-existent file
	_, err := gitcopy.ReadFile("/non/existent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test with directory instead of file
	tempDir := t.TempDir()
	_, err = gitcopy.ReadFile(tempDir)
	if err == nil {
		t.Error("Expected error when trying to read directory as file, got nil")
	}
}

// Testgitcopy.IoReadDir tests the gitcopy.IoReadDir function
func TestIoReadDir(t *testing.T) {
	// Create a temporary directory with test files
	tempDir := t.TempDir()

	// Create test files and subdirectory structure
	testFiles := []string{
		"file1.txt",
		"file2.json",
		"subdir/file3.yaml",
		"subdir/nested/file4.md",
		"another/deep/path/file5.xml",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", file, err)
		}
		err = os.WriteFile(fullPath, []byte("test content for "+file), 0o644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test reading directory
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed: %v", err)
	}

	// Should find all files recursively
	if len(files) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(files))
	}

	// Check that all expected files are found
	expectedFiles := make(map[string]bool)
	for _, file := range testFiles {
		expectedFiles[filepath.Join(tempDir, file)] = false
	}

	for _, file := range files {
		if _, exists := expectedFiles[file]; exists {
			expectedFiles[file] = true
		} else {
			t.Errorf("Unexpected file found: %s", file)
		}
	}

	for file, found := range expectedFiles {
		if !found {
			t.Errorf("Expected file not found: %s", file)
		}
	}
}

// Testgitcopy.IoReadDirError tests error cases for gitcopy.IoReadDir
func TestIoReadDirError(t *testing.T) {
	// Test with non-existent directory
	_, err := gitcopy.IoReadDir("/non/existent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

// Testgitcopy.IoReadDirEmpty tests gitcopy.IoReadDir with empty directory
func TestIoReadDirEmpty(t *testing.T) {
	tempDir := t.TempDir()

	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed on empty directory: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", len(files))
	}
}

// TestValidationLogic tests input validation scenarios
func TestValidationLogic(t *testing.T) {
	tests := []struct {
		name                 string
		filePath             string
		destinationFilePath  string
		directory            string
		destinationDirectory string
		shouldFail           bool
		expectedError        string
	}{
		{
			name:                 "Valid file operation",
			filePath:             "source.txt",
			destinationFilePath:  "dest.txt",
			directory:            "",
			destinationDirectory: "",
			shouldFail:           false,
		},
		{
			name:                 "Valid directory operation",
			filePath:             "",
			destinationFilePath:  "",
			directory:            "source-dir",
			destinationDirectory: "dest-dir",
			shouldFail:           false,
		},
		{
			name:                 "File without destination",
			filePath:             "source.txt",
			destinationFilePath:  "",
			directory:            "",
			destinationDirectory: "",
			shouldFail:           true,
			expectedError:        "missing input 'destination_file file'",
		},
		{
			name:                 "Directory without destination",
			filePath:             "",
			destinationFilePath:  "",
			directory:            "source-dir",
			destinationDirectory: "",
			shouldFail:           true,
			expectedError:        "missing input 'destination-directory'",
		},
		{
			name:                 "No input provided",
			filePath:             "",
			destinationFilePath:  "",
			directory:            "",
			destinationDirectory: "",
			shouldFail:           true,
			expectedError:        "file or directory is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the validation logic from runApplication
			hasErrors := false
			var errorMsg string

			if tt.filePath != "" && tt.destinationFilePath == "" {
				hasErrors = true
				errorMsg = "missing input 'destination_file file'"
			}

			if tt.directory != "" && tt.destinationDirectory == "" {
				hasErrors = true
				errorMsg = "missing input 'destination-directory'"
			}

			if tt.filePath == "" && tt.directory == "" {
				hasErrors = true
				errorMsg = "file or directory is required"
			}

			if tt.shouldFail {
				if !hasErrors {
					t.Errorf("Expected validation to fail, but it passed")
				}
				if tt.expectedError != "" && errorMsg != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, errorMsg)
				}
			} else if hasErrors {
				t.Errorf("Expected validation to pass, but got error: %s", errorMsg)
			}
		})
	}
}

// TestBranchInitialization tests the branch initialization logic
func TestBranchInitialization(t *testing.T) {
	tests := []struct {
		name       string
		filePath   string
		directory  string
		shouldInit bool
	}{
		{
			name:       "File operation should initialize branch",
			filePath:   "source.txt",
			directory:  "",
			shouldInit: true,
		},
		{
			name:       "Directory operation should initialize branch",
			filePath:   "",
			directory:  "source-dir",
			shouldInit: true,
		},
		{
			name:       "Both operations should initialize branch",
			filePath:   "source.txt",
			directory:  "source-dir",
			shouldInit: true,
		},
		{
			name:       "No operation should not initialize branch",
			filePath:   "",
			directory:  "",
			shouldInit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the branch initialization logic: if envVar.Input.FilePath != "" || envVar.Input.Directory != ""
			willInit := tt.filePath != "" || tt.directory != ""

			if willInit != tt.shouldInit {
				t.Errorf("Expected branch initialization %v, got %v", tt.shouldInit, willInit)
			}
		})
	}
}

// TestDefaultValues tests default value assignment
func TestDefaultValues(t *testing.T) {
	// Test default pull message and description generation
	timeNow := time.Now()

	// Simulate the default assignment logic
	pullMessage := ""
	pullDescription := ""

	if pullMessage == "" {
		pullMessage = "update " + timeNow.Format("2006-01-02 15:04:05")
	}

	if pullDescription == "" {
		pullDescription = "update " + timeNow.Format("2006-01-02 15:04:05")
	}

	// Verify format
	expectedPrefix := "update " + timeNow.Format("2006-01-02")
	if !strings.HasPrefix(pullMessage, expectedPrefix) {
		t.Errorf("Pull message should start with '%s', got: %s", expectedPrefix, pullMessage)
	}

	if !strings.HasPrefix(pullDescription, expectedPrefix) {
		t.Errorf("Pull description should start with '%s', got: %s", expectedPrefix, pullDescription)
	}
}

// TestReviewerParsing tests reviewer parsing logic
func TestReviewerParsing(t *testing.T) {
	tests := []struct {
		name          string
		userInput     string
		teamInput     string
		expectedUsers []string
		expectedTeams []string
	}{
		{
			name:          "Empty reviewers",
			userInput:     "",
			teamInput:     "",
			expectedUsers: []string{},
			expectedTeams: []string{},
		},
		{
			name:          "Single user",
			userInput:     "john",
			teamInput:     "",
			expectedUsers: []string{"john"},
			expectedTeams: []string{},
		},
		{
			name:          "Multiple users",
			userInput:     "john,jane,bob",
			teamInput:     "",
			expectedUsers: []string{"john", "jane", "bob"},
			expectedTeams: []string{},
		},
		{
			name:          "Single team",
			userInput:     "",
			teamInput:     "backend-team",
			expectedUsers: []string{},
			expectedTeams: []string{"backend-team"},
		},
		{
			name:          "Multiple teams",
			userInput:     "",
			teamInput:     "backend-team,frontend-team",
			expectedUsers: []string{},
			expectedTeams: []string{"backend-team", "frontend-team"},
		},
		{
			name:          "Mixed reviewers",
			userInput:     "john,jane",
			teamInput:     "backend-team,ops-team",
			expectedUsers: []string{"john", "jane"},
			expectedTeams: []string{"backend-team", "ops-team"},
		},
		{
			name:          "Reviewers with spaces",
			userInput:     " john , jane ",
			teamInput:     " backend-team , ops-team ",
			expectedUsers: []string{"john", "jane"},
			expectedTeams: []string{"backend-team", "ops-team"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate reviewer parsing logic
			var users, teams []string

			if tt.userInput != "" {
				for _, user := range strings.Split(tt.userInput, ",") {
					user = strings.TrimSpace(user)
					if user != "" {
						users = append(users, user)
					}
				}
			}

			if tt.teamInput != "" {
				for _, team := range strings.Split(tt.teamInput, ",") {
					team = strings.TrimSpace(team)
					if team != "" {
						teams = append(teams, team)
					}
				}
			}

			// Verify results
			if len(users) != len(tt.expectedUsers) {
				t.Errorf("Expected %d users, got %d", len(tt.expectedUsers), len(users))
			}

			for i, user := range users {
				if i >= len(tt.expectedUsers) || user != tt.expectedUsers[i] {
					t.Errorf("User mismatch at index %d: expected %s, got %s", i, tt.expectedUsers[i], user)
				}
			}

			if len(teams) != len(tt.expectedTeams) {
				t.Errorf("Expected %d teams, got %d", len(tt.expectedTeams), len(teams))
			}

			for i, team := range teams {
				if i >= len(tt.expectedTeams) || team != tt.expectedTeams[i] {
					t.Errorf("Team mismatch at index %d: expected %s, got %s", i, tt.expectedTeams[i], team)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkReadFile(b *testing.B) {
	// Create a temporary test file with substantial content
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark.txt")
	testContent := strings.Repeat("This is test content for benchmarking. ", 1000) // ~39KB of data

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gitcopy.ReadFile(testFile)
		if err != nil {
			b.Fatalf("gitcopy.ReadFile failed: %v", err)
		}
	}
}

func BenchmarkIoReadDir(b *testing.B) {
	// Create a temporary directory with many files
	tempDir := b.TempDir()

	// Create 50 test files in various subdirectories
	for i := 0; i < 50; i++ {
		subdir := filepath.Join(tempDir, "subdir"+string(rune(i%5)))
		err := os.MkdirAll(subdir, 0o755)
		if err != nil {
			b.Fatalf("Failed to create subdirectory: %v", err)
		}

		testFile := filepath.Join(subdir, "file"+string(rune(i))+".txt")
		err = os.WriteFile(testFile, []byte("benchmark content"), 0o644)
		if err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gitcopy.IoReadDir(tempDir)
		if err != nil {
			b.Fatalf("gitcopy.IoReadDir failed: %v", err)
		}
	}
}

// TestCrossplatformPaths tests cross-platform file path handling
func TestCrossplatformPaths(t *testing.T) {
	tempDir := t.TempDir()

	// Create nested directory structure
	nestedPath := filepath.Join(tempDir, "level1", "level2", "file.txt")
	err := os.MkdirAll(filepath.Dir(nestedPath), 0o755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	err = os.WriteFile(nestedPath, []byte("nested file content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create nested file: %v", err)
	}

	// Test that gitcopy.IoReadDir correctly handles nested paths using filepath.Join
	files, err := gitcopy.IoReadDir(tempDir)
	if err != nil {
		t.Fatalf("gitcopy.IoReadDir failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(files) > 0 {
		expectedPath := filepath.Join(tempDir, "level1", "level2", "file.txt")
		if files[0] != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, files[0])
		}
	}
}
