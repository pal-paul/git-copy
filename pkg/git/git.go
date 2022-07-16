package git

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Git struct {
	owner string
	repo  string
	token string
}

func New(owner string, repo string, token string) *Git {
	return &Git{
		owner: owner,
		repo:  repo,
		token: token,
	}
}

func (g *Git) GetBranch(branch string) (*BranchInfo, error) {
	var branchInfo BranchInfo
	resp, err := g.get("repos", fmt.Sprintf("%s/%s/git/refs/heads/%s", g.owner, g.repo, branch), nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get branch %s: %s", branch, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &branchInfo)
	if err != nil {
		return nil, err
	}
	return &branchInfo, nil
}

func (g *Git) CreateBranch(branch string, sha string) (*BranchInfo, error) {
	reqBody := map[string]string{
		"ref": fmt.Sprintf("refs/heads/%s", branch),
		"sha": sha,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := g.post("repos", fmt.Sprintf("%s/%s/git/refs", g.owner, g.repo), nil, reqBodyJson)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("failed to create a branch %s: %s", branch, resp.Status)
	}
	return nil, nil
}

func (g *Git) GetAFile(branch string, filePath string) (*FileInfo, error) {
	var fileInfo FileInfo
	qs := url.Values{}
	qs.Add("ref", branch)
	resp, err := g.get("repos", fmt.Sprintf("%s/%s/contents/%s", g.owner, g.repo, filePath), qs)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get file %s: %s", filePath, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fileInfo)
	if err != nil {
		return nil, err
	}
	return &fileInfo, nil
}

func (g *Git) CreateUpdateAFile(branch string, filePath string, content []byte, message string, sha string) (*FileResponse, error) {
	var fileResponse FileResponse
	b64content := b64.StdEncoding.EncodeToString(content)
	reqBody := map[string]string{
		"message": message,
		"content": b64content,
		"branch":  branch,
	}
	if sha != "" {
		reqBody["sha"] = sha
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := g.put("repos", fmt.Sprintf("%s/%s/contents/%s", g.owner, g.repo, filePath), nil, reqBodyJson)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 201 {
		return nil, fmt.Errorf("failed to update file %s: %s", filePath, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return nil, err
	}
	return &fileResponse, nil
}

func (g *Git) CreatePullRequest(baseBranch string, branch string, title string, description string) error {
	reqBody := map[string]any{
		"title":                 title,
		"body":                  description,
		"head":                  branch,
		"base":                  baseBranch,
		"maintainer_can_modify": true,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	resp, err := g.post("repos", fmt.Sprintf("%s/%s/pulls", g.owner, g.repo), nil, reqBodyJson)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create pull request: %s", resp.Status)
	}
	return nil
}
