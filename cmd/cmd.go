package main

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pal-paul/git-copy/pkg/git"
)

var (
	owner string
	repo  string
	token string
)

func main() {
	var (
		file                 string
		destinationFile      string
		directory            string
		destinationDirectory string
		pullMessage          string
		pullDescription      string
		reviewers            string
		teamReviewers        string
		refPullBranch        string
		refBranch            string
		branch               string
	)

	refPullBranch = os.Getenv("INPUT_REF_BRANCH")
	if refPullBranch == "" {
		refPullBranch = "master"
	}

	owner = os.Getenv("INPUT_OWNER")
	repo = os.Getenv("INPUT_REPO")
	token = os.Getenv("INPUT_TOKEN")

	file = os.Getenv("INPUT_FILE_PATH")
	destinationFile = os.Getenv("INPUT_DESTINATION_FILE_PATH")

	directory = os.Getenv("INPUT_DIRECTORY")
	destinationDirectory = os.Getenv("INPUT_DESTINATION_DIRECTORY")

	pullMessage = os.Getenv("INPUT_PULL_MESSAGE")
	pullDescription = os.Getenv("INPUT_PULL_DESCRIPTION")

	reviewers = os.Getenv("INPUT_REVIEWERS")
	teamReviewers = os.Getenv("INPUT_TEAM_REVIEWERS")

	if owner == "" {
		log.Fatal("owner is required")
		return
	}
	if repo == "" {
		log.Fatal("missing input 'repo'")
		return
	}
	if token == "" {
		log.Fatal("missing input 'token'")
		return
	}

	if file != "" && destinationFile == "" {
		log.Fatal("missing input 'destinationfile file'")
		return
	}

	if directory != "" && destinationDirectory == "" {
		log.Fatal("missing input 'destination-directory'")
		return
	}

	if file == "" && directory == "" {
		log.Fatal("file or directory is required")
		return
	}

	if pullMessage == "" {
		pullMessage = fmt.Sprintf("update %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	if pullDescription == "" {
		pullDescription = fmt.Sprintf("update %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	var gitReviewers git.Reviewers
	if reviewers != "" {
		reviewers = strings.ReplaceAll(reviewers, " ", "")
		gitReviewers.Users = append(gitReviewers.Users, strings.Split(reviewers, ",")...)
	}
	if teamReviewers != "" {
		teamReviewers = strings.ReplaceAll(teamReviewers, " ", "")
		gitReviewers.Teams = append(gitReviewers.Teams, strings.Split(teamReviewers, ",")...)
	}

	branch = os.Getenv("INPUT_BRANCH")
	if branch == "" {
		branch = uuid.New().String()
	}

	// DO NOT EDIT BELOW THIS LINE
	gitObj := git.New(owner, repo, token)
	if file != "" || directory != "" {
		refBranch = refPullBranch
		refDefaultBranch, err := gitObj.GetBranch(refBranch)
		if err != nil {
			log.Fatal(err)
		}
		copyToBranch, err := gitObj.GetBranch(branch)
		if err != nil {
			log.Fatal(err)
		}
		if copyToBranch == nil {
			_, err = gitObj.CreateBranch(branch, refDefaultBranch.Object.Sha)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			refBranch = branch
		}
	}

	if file != "" {
		err := uploadFile(refBranch, branch, gitObj, file, destinationFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	if directory != "" {
		files, err := ioReadDir(directory)
		if err != nil {
			log.Fatal(err)
		}
		err = uploadFiles(refBranch, branch, gitObj, files, directory, destinationDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}
	if file != "" || directory != "" {
		if refBranch == refPullBranch {
			prNumber, err := gitObj.CreatePullRequest(
				refPullBranch,
				branch,
				pullMessage,
				pullDescription,
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
	} else {
		log.Fatal("nothing to do")
	}
}

func ioReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			nestedFiles, err := ioReadDir(root + "/" + file.Name())
			if err != nil {
				return files, err
			}
			files = append(files, nestedFiles...)
		} else {
			files = append(files, root+"/"+file.Name())
		}
	}
	return files, nil
}

func uploadFile(
	refBranch string,
	branch string,
	gitObj *git.Git,
	file string,
	detinationfile string,
) error {
	fileContent, err := readFile(file)
	if err != nil {
		return err
	}
	fileObj, err := gitObj.GetAFile(refBranch, detinationfile)
	if err != nil {
		return err
	}
	if fileObj == nil {
		fmt.Println("creating file:", file)
		_, err = gitObj.CreateUpdateAFile(
			branch,
			detinationfile,
			fileContent,
			fmt.Sprintf("%s file created", detinationfile),
			"",
		)
		if err != nil {
			return err
		}
	} else {
		if string(b64.StdEncoding.EncodeToString(fileContent)) != string(fileObj.Content) {
			fmt.Println("updating file:", file)
			_, err = gitObj.CreateUpdateAFile(branch, detinationfile, fileContent, fmt.Sprintf("%s file updated", detinationfile), fileObj.Sha)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("no changes file:", file)
		}
	}
	return nil
}

func uploadFiles(refBranch string, branch string, gitObj *git.Git, files []string, srcDir string, destDir string) error {
	batch := git.BatchFileUpdate{
		Branch:  branch,
		Message: fmt.Sprintf("Batch update files in %s", destDir),
		Files:   make([]git.FileOperation, 0),
	}

	for _, file := range files {
		destinationFile := strings.Replace(file, srcDir, destDir, 1)

		fileContent, err := readFile(file)
		if err != nil {
			return err
		}

		fileObj, err := gitObj.GetAFile(refBranch, destinationFile)
		if err != nil {
			return err
		}

		encoded := b64.StdEncoding.EncodeToString(fileContent)

		fileOp := git.FileOperation{
			Path:    destinationFile,
			Content: encoded,
		}

		if fileObj != nil {
			if encoded != string(fileObj.Content) {
				fileOp.Sha = fileObj.Sha
				batch.Files = append(batch.Files, fileOp)
				fmt.Println("updating file:", file)
			} else {
				fmt.Println("no changes file:", file)
			}
		} else {
			batch.Files = append(batch.Files, fileOp)
			fmt.Println("creating file:", file)
		}
	}

	if len(batch.Files) > 0 {
		return gitObj.CreateUpdateMultipleFiles(batch)
	}

	return nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}
