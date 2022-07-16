package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pal-paul/git-copy/pkg/git"
)

var (
	owner string
	repo  string
	token string

	refbranch string
	branch    string
)

func main() {
	var (
		file                string
		detinationfile      string
		directory           string
		detinationdirectory string
	)
	refbranch = "master"
	owner = os.Getenv("INPUT_OWNER")
	repo = os.Getenv("INPUT_REPO")
	token = os.Getenv("INPUT_TOKEN")

	file = os.Getenv("INPUT_FILE_PATH")
	detinationfile = os.Getenv("INPUT_DETINATION_FILE_PATH")

	directory = os.Getenv("INPUT_DIRECTORY")
	detinationdirectory = os.Getenv("INPUT_DETINATION_DIRECTORY")

	if owner == "" {
		fmt.Println("missing input 'owner'")
	}
	if repo == "" {
		fmt.Println("missing input 'repo'")
	}
	if token == "" {
		fmt.Println("missing input 'token'")
	}

	if file == "" {
		fmt.Println("missing input 'file'")
	}
	if detinationfile == "" {
		fmt.Println("missing input 'detinationfile file'")
	}

	if directory == "" {
		fmt.Println("missing input 'directory'")
	}
	if detinationdirectory == "" {
		fmt.Println("missing input 'detination-directory'")
	}

	branch = uuid.New().String()

	// DO NOT EDIT BELOW THIS LINE
	gitobj := git.New(owner, repo, token)
	if file != "" {
		refBranch, err := gitobj.GetBranch(refbranch)
		if err != nil {
			fmt.Println(err)
		}
		_, err = gitobj.CreateBranch(branch, refBranch.Object.Sha)
		if err != nil {
			fmt.Println(err)
		}
		err = uploadFile(gitobj, file, detinationfile)
		if err != nil {
			fmt.Println(err)
		}
	}
	if directory != "" {
		files, err := ioReadDir(directory)
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range files {
			detinationfile := strings.Replace(file, directory, detinationdirectory, 1)
			err = uploadFile(gitobj, file, detinationfile)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	err := gitobj.CreatePullRequest(refbranch, branch, "updates", "just a test")
	if err != nil {
		fmt.Println(err)
	}
}

func ioReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			nestedfiles, err := ioReadDir(root + "/" + file.Name())
			if err != nil {
				return files, err
			}
			files = append(files, nestedfiles...)
		} else {
			files = append(files, root+"/"+file.Name())
		}
	}
	return files, nil
}

func uploadFile(gitobj *git.Git, file string, detinationfile string) error {
	fmt.Println("updating file:", file)
	filecontent, err := readfile(file)
	if err != nil {
		return err
	}
	fileObj, err := gitobj.GetAFile(refbranch, detinationfile)
	if err != nil {
		return err
	}

	if fileObj == nil {
		_, err = gitobj.CreateUpdateAFile(branch, detinationfile, filecontent, fmt.Sprintf("%s file created", detinationfile), "")
		if err != nil {
			return err
		}
	} else {
		_, err = gitobj.CreateUpdateAFile(branch, detinationfile, filecontent, fmt.Sprintf("%s file updated", detinationfile), fileObj.Sha)
		if err != nil {
			return err
		}
	}
	return nil
}

func readfile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}
