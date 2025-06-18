package gitcopy

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pal-paul/go-libraries/pkg/env"
	"github.com/pal-paul/go-libraries/pkg/git"
)

type Environment struct {
	GitHub struct {
		Token    string `env:"GITHUB_TOKEN,required=true"`
		Api      string `env:"GITHUB_API_URL,required=true"`
		Repo     string `env:"GITHUB_REPOSITORY,required=true"`
		Workflow string `env:"GITHUB_WORKFLOW,required=true"`
		Branch   string `env:"GITHUB_REF,required=true"`
		Commit   string `env:"GITHUB_SHA,required=true"`
		RunId    string `env:"GITHUB_RUN_ID,required=true"`
		JobName  string `env:"GITHUB_JOB,required=true"`
		Server   string `env:"GITHUB_SERVER_URL,required=true"`
	}
	Input struct {
		Owner                string `env:"INPUT_OWNER,required=true"`
		Repo                 string `env:"INPUT_REPO,required=true"`
		FilePath             string `env:"INPUT_FILE_PATH,required=false"`
		DestinationFilePath  string `env:"INPUT_DESTINATION_FILE_PATH,required=false"`
		Directory            string `env:"INPUT_DIRECTORY,required=false"`
		DestinationDirectory string `env:"INPUT_DESTINATION_DIRECTORY,required=false"`
		PullMessage          string `env:"INPUT_PULL_MESSAGE,required=false"`
		PullDescription      string `env:"INPUT_PULL_DESCRIPTION,required=false"`
		Reviewers            string `env:"INPUT_REVIEWERS,required=false"`
		TeamReviewers        string `env:"INPUT_TEAM_REVIEWERS,required=false"`
		RefBranch            string `env:"INPUT_REF_BRANCH,default=master"`
		Branch               string `env:"INPUT_BRANCH,default=update-branch"`
	}
}

var envVar Environment

// InitializeEnvironment loads environment variables for the application
func InitializeEnvironment() error {
	_, err := env.Unmarshal(&envVar)
	return err
}

// GetEnvironment returns the current environment configuration
func GetEnvironment() Environment {
	return envVar
}

// SetEnvironment sets the environment configuration (useful for testing)
func SetEnvironment(env Environment) {
	envVar = env
}

// IoReadDir recursively reads all files in a directory and its subdirectories
func IoReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			nestedFiles, err := IoReadDir(filepath.Join(root, file.Name()))
			if err != nil {
				// Log error but continue processing other files
				log.Printf("Error reading directory %s: %v", filepath.Join(root, file.Name()), err)
				continue
			}
			files = append(files, nestedFiles...)
		} else {
			files = append(files, filepath.Join(root, file.Name()))
		}
	}
	return files, nil
}

// ReadFile reads the contents of a file
func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

// RunApplication executes the master application logic
func RunApplication() {
	var refBranch string

	if envVar.Input.FilePath != "" && envVar.Input.DestinationFilePath == "" {
		log.Fatal("missing input 'destination_file file'")
		return
	}

	if envVar.Input.Directory != "" && envVar.Input.DestinationDirectory == "" {
		log.Fatal("missing input 'destination-directory'")
		return
	}

	if envVar.Input.FilePath == "" && envVar.Input.Directory == "" {
		log.Fatal("file or directory is required")
		return
	}

	if envVar.Input.PullMessage == "" {
		envVar.Input.PullMessage = fmt.Sprintf("update %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	if envVar.Input.PullDescription == "" {
		envVar.Input.PullDescription = fmt.Sprintf("update %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	var gitReviewers git.Reviewers
	if envVar.Input.Reviewers != "" {
		envVar.Input.Reviewers = strings.ReplaceAll(envVar.Input.Reviewers, " ", "")
		gitReviewers.Users = append(gitReviewers.Users, strings.Split(envVar.Input.Reviewers, ",")...)
	}
	if envVar.Input.TeamReviewers != "" {
		envVar.Input.TeamReviewers = strings.ReplaceAll(envVar.Input.TeamReviewers, " ", "")
		gitReviewers.Teams = append(gitReviewers.Teams, strings.Split(envVar.Input.TeamReviewers, ",")...)
	}

	if envVar.Input.Branch == "" {
		envVar.Input.Branch = uuid.New().String()
	}

	// DO NOT EDIT BELOW THIS LINE
	gitObj := git.New(
		git.WithOwner(envVar.Input.Owner),
		git.WithRepo(envVar.Input.Repo),
		git.WithToken(envVar.GitHub.Token),
	)

	// Initialize branch handling for both file and directory operations
	if envVar.Input.FilePath != "" || envVar.Input.Directory != "" {
		refBranch = envVar.Input.RefBranch
		refDefaultBranch, err := gitObj.GetBranch(refBranch)
		if err != nil {
			log.Fatal(err)
		}
		copyToBranch, err := gitObj.GetBranch(envVar.Input.Branch)
		if err != nil {
			log.Fatal(err)
		}
		if copyToBranch == nil {
			_, err = gitObj.CreateBranch(envVar.Input.Branch, refDefaultBranch.Object.Sha)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			refBranch = envVar.Input.Branch
		}
	}

	if envVar.Input.FilePath != "" {
		fileContent, err := ReadFile(envVar.Input.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		fileObj, err := gitObj.GetAFile(refBranch, envVar.Input.DestinationFilePath)
		if err != nil {
			log.Fatal(err)
		}
		if fileObj == nil {
			fmt.Println("creating file:", envVar.Input.FilePath)
			_, err = gitObj.CreateUpdateAFile(
				envVar.Input.Branch,
				envVar.Input.DestinationFilePath,
				fileContent,
				fmt.Sprintf("%s file created", envVar.Input.DestinationFilePath),
				"",
			)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			encoded := b64.StdEncoding.EncodeToString(fileContent)
			if encoded != fileObj.Content {
				fmt.Println("updating file:", envVar.Input.FilePath)
				_, err = gitObj.CreateUpdateAFile(
					envVar.Input.Branch,
					envVar.Input.DestinationFilePath,
					fileContent,
					fmt.Sprintf("%s file updated", envVar.Input.DestinationFilePath),
					fileObj.Sha,
				)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("INFO: No changes detected for %s", envVar.Input.FilePath)
			}
		}
	}

	if envVar.Input.Directory != "" {
		files, err := IoReadDir(envVar.Input.Directory)
		if err != nil {
			log.Fatal(err)
		}
		batch := git.BatchFileUpdate{
			Branch:  envVar.Input.Branch,
			Message: fmt.Sprintf("Batch update files in %s", envVar.Input.DestinationDirectory),
			Files:   make([]git.FileOperation, 0),
		}

		for _, file := range files {
			relativePath, err := filepath.Rel(envVar.Input.Directory, file)
			if err != nil {
				log.Printf("ERROR: could not get relative path for %s: %v", file, err)
				continue
			}
			destinationFile := filepath.Join(envVar.Input.DestinationDirectory, relativePath)

			fileContent, err := ReadFile(file)
			if err != nil {
				log.Printf("ERROR: could not read file %s: %v", file, err)
				continue
			}

			fileObj, err := gitObj.GetAFile(refBranch, destinationFile)
			if err != nil {
				// log.Printf("ERROR: could not get file %s: %v", destinationFile, err)
				continue
			}

			encoded := b64.StdEncoding.EncodeToString(fileContent)

			fileOp := git.FileOperation{
				Path:    destinationFile,
				Content: encoded,
			}
			if fileObj != nil {
				if encoded != fileObj.Content {
					fileOp.Sha = fileObj.Sha
					batch.Files = append(batch.Files, fileOp)
				}
			} else {
				batch.Files = append(batch.Files, fileOp)
			}
		}

		if len(batch.Files) > 0 {
			fmt.Printf("Processing batch update for %d files\n", len(batch.Files))
			err = gitObj.CreateUpdateMultipleFiles(batch)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("INFO: No files need updating in %s", envVar.Input.DestinationDirectory)
		}
	}

	if envVar.Input.FilePath != "" || envVar.Input.Directory != "" {
		if refBranch != envVar.Input.Branch {
			prNumber, err := gitObj.CreatePullRequest(
				refBranch,
				envVar.Input.Branch,
				envVar.Input.PullMessage,
				envVar.Input.PullDescription,
			)
			if err != nil {
				log.Fatal(err)
			}
			if gitReviewers.Users != nil || gitReviewers.Teams != nil {
				err = gitObj.AddReviewers(prNumber, gitReviewers)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
