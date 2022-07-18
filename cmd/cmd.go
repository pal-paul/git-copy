package main

import (
	"fmt"
	"io/ioutil"
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

	refbranch string
	branch    string
)

func main() {
	var (
		file                string
		detinationfile      string
		directory           string
		detinationdirectory string
		pullmessage         string
		pulldescption       string
	)
	refbranch = "master"
	owner = os.Getenv("INPUT_OWNER")
	repo = os.Getenv("INPUT_REPO")
	token = os.Getenv("INPUT_TOKEN")

	file = os.Getenv("INPUT_FILE_PATH")
	detinationfile = os.Getenv("INPUT_DETINATION_FILE_PATH")

	directory = os.Getenv("INPUT_DIRECTORY")
	detinationdirectory = os.Getenv("INPUT_DETINATION_DIRECTORY")

	pullmessage = os.Getenv("INPUT_PULL_MESSAGE")
	pulldescption = os.Getenv("INPUT_PULL_DESCRIPTION")

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

	if file != "" && detinationfile == "" {
		log.Fatal("missing input 'detinationfile file'")
		return
	}

	if directory != "" && detinationdirectory == "" {
		log.Fatal("missing input 'detination-directory'")
		return
	}

	if file == "" && directory == "" {
		log.Fatal("file or directory is required")
		return
	}

	if pullmessage == "" {
		pullmessage = fmt.Sprintf("update %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	if pulldescption == "" {
		pulldescption = fmt.Sprintf("update %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	branch = uuid.New().String()

	// DO NOT EDIT BELOW THIS LINE
	gitobj := git.New(owner, repo, token)
	if file != "" || directory != "" {
		refBranch, err := gitobj.GetBranch(refbranch)
		if err != nil {
			log.Fatal(err)
		}
		_, err = gitobj.CreateBranch(branch, refBranch.Object.Sha)
		if err != nil {
			log.Fatal(err)
		}
	}

	if file != "" {
		err := uploadFile(gitobj, file, detinationfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	if directory != "" {
		files, err := ioReadDir(directory)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			detinationfile := strings.Replace(file, directory, detinationdirectory, 1)
			err = uploadFile(gitobj, file, detinationfile)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if file != "" || directory != "" {
		err := gitobj.CreatePullRequest(refbranch, branch, pullmessage, pulldescption)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("nothing to do")
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
	filecontent, err := readfile(file)
	if err != nil {
		return err
	}
	fileObj, err := gitobj.GetAFile(refbranch, detinationfile)
	if err != nil {
		return err
	}
	if fileObj == nil {
		fmt.Println("creating file:", file)
		_, err = gitobj.CreateUpdateAFile(branch, detinationfile, filecontent, fmt.Sprintf("%s file created", detinationfile), "")
		if err != nil {
			return err
		}
	} else {
		if string(filecontent) != string(fileObj.Content) {
			fmt.Println("updating file:", file)
			_, err = gitobj.CreateUpdateAFile(branch, detinationfile, filecontent, fmt.Sprintf("%s file updated", detinationfile), fileObj.Sha)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("no changes file:", file)
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
