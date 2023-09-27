package git

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	resp, err := g.get(
		"repos",
		fmt.Sprintf("%s/%s/git/refs/heads/%s", g.owner, g.repo, branch),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get branch %s: %s", branch, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &branchInfo)
	if err != nil {
		fmt.Println(string(body))
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fileInfo)
	if err != nil {
		return nil, err
	}
	return &fileInfo, nil
}

func (g *Git) CreateUpdateAFile(
	branch string,
	filePath string,
	content []byte,
	message string,
	sha string,
) (*FileResponse, error) {
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
	resp, err := g.put(
		"repos",
		fmt.Sprintf("%s/%s/contents/%s", g.owner, g.repo, filePath),
		nil,
		reqBodyJson,
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 201 {
		return nil, fmt.Errorf("failed to update file %s: %s", filePath, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &fileResponse)
	if err != nil {
		return nil, err
	}
	return &fileResponse, nil
}

// CreatePullRequest creates a pull request and returns the pull request number.
func (g *Git) CreatePullRequest(
	baseBranch string,
	branch string,
	title string,
	description string,
) (int, error) {
	reqBody := map[string]any{
		"title":                 title,
		"body":                  description,
		"head":                  branch,
		"base":                  baseBranch,
		"maintainer_can_modify": true,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}
	resp, err := g.post("repos", fmt.Sprintf("%s/%s/pulls", g.owner, g.repo), nil, reqBodyJson)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 201 {
		return 0, fmt.Errorf("failed to create pull request: %s", resp.Status)
	}
	var pullResponse PullResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(body, &pullResponse)
	if err != nil {
		return 0, err
	}
	return pullResponse.Number, nil
}

func (g *Git) AddReviewers(number int, prReviewers Reviewers) error {
	reqBody := make(map[string][]string)
	if len(prReviewers.Users) > 0 {
		reqBody["reviewers"] = prReviewers.Users
	}
	if len(prReviewers.Teams) > 0 {
		reqBody["team_reviewers"] = prReviewers.Teams
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	resp, err := g.post(
		"repos",
		fmt.Sprintf("%s/%s/pulls/%d/requested_reviewers", g.owner, g.repo, number),
		nil,
		reqBodyJson,
	)
	if err != nil {
		return err
	}
	if resp.StatusCode == 422 {
		fmt.Println("reviewer are not collaborator")
		return nil
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to add reviewers: %s", resp.Status)
	}
	return nil
}
