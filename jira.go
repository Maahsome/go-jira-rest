package jira

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Jira struct {
	BaseUrl   string
	ApiPath   string
	AgilePath string
	User      string
	Token     string
	Client    *resty.Client
}

type Boards struct {
	MaxResults int           `json:"maxResults"`
	StartAt    int           `json:"startAt"`
	Total      int           `json:"total"`
	IsLast     bool          `json:"isLast"`
	Values     []BoardValues `json:"values"`
}

type Location struct {
	ProjectID      int    `json:"projectId"`
	DisplayName    string `json:"displayName"`
	ProjectName    string `json:"projectName"`
	ProjectKey     string `json:"projectKey"`
	ProjectTypeKey string `json:"projectTypeKey"`
	AvatarURI      string `json:"avatarURI"`
	Name           string `json:"name"`
}

type BoardValues struct {
	ID       int      `json:"id"`
	Self     string   `json:"self"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location Location `json:"location"`
}

type SprintIssues struct {
	Expand     string     `json:"expand"`
	StartAt    int        `json:"startAt"`
	MaxResults int        `json:"maxResults"`
	Total      int        `json:"total"`
	Issues     []struct{} `json:"issues"`
}

// New generate a new jira client
func New(baseUrl, apiPath, agilePath, user, token string) *Jira {

	restClient := resty.New()

	if apiPath == "" {
		apiPath = "/rest/api/3"
	}
	if agilePath == "" {
		agilePath = "/rest/agile/1.0"
	}
	restClient.SetBasicAuth(user, token)

	return &Jira{
		BaseUrl:   baseUrl,
		ApiPath:   apiPath,
		AgilePath: agilePath,
		Token:     token,
		Client:    restClient,
	}
}

func (r *Jira) GetIssue(issue string) (string, error) {

	fetchUri := fmt.Sprintf("%s%s/issue/%s", r.BaseUrl, r.ApiPath, issue)
	// logrus.Warn(fetchUri)
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		Get(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

func (r *Jira) AddComment(issue string, comment string) (string, error) {

	fetchUri := fmt.Sprintf("%s%s/issue/%s/comment", r.BaseUrl, r.ApiPath, issue)
	// logrus.Warn(fetchUri)
	body := fmt.Sprintf("{ \"body\": \"%s\" }", comment)
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

func (r *Jira) AssignIssue(issue string, account string) error {

	var body string
	fetchUri := fmt.Sprintf("%s%s/issue/%s/assignee", r.BaseUrl, r.ApiPath, issue)
	// logrus.Warn(fetchUri)
	if account == "null" {
		body = fmt.Sprintf("{ \"accountId\": %s }", account)
	} else {
		body = fmt.Sprintf("{ \"accountId\": \"%s\" }", account)
	}
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops", resp.Body()[:])
		return resperr
	}

	return nil
}

func (r *Jira) GetAccount(search string) (string, error) {

	fetchUri := fmt.Sprintf("%s%s/user/search?query=%s", r.BaseUrl, r.ApiPath, search)
	// logrus.Warn(fetchUri)
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		Get(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops", resp.Body()[:])
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

func (r *Jira) GetTransitions(issue string) (string, error) {

	fetchUri := fmt.Sprintf("%s%s/issue/%s/transitions", r.BaseUrl, r.ApiPath, issue)
	// logrus.Warn(fetchUri)
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		Get(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

func (r *Jira) TransitionIssue(issue string, transitionId string) (string, error) {

	fetchUri := fmt.Sprintf("%s%s/issue/%s/transitions", r.BaseUrl, r.ApiPath, issue)
	// logrus.Warn(fetchUri)
	body := fmt.Sprintf("{ \"transition\": { \"id\": %s } }", transitionId)
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

func (r *Jira) GetBoards() (string, error) {

	startAt := 0
	var returnValues []BoardValues

	for {
		fetchUri := fmt.Sprintf("%s%s/board?startAt=%d", r.BaseUrl, r.AgilePath, startAt)

		resp, resperr := r.Client.R().
			SetHeader("Content-Type", "application/json").
			Get(fetchUri)

		if resperr != nil {
			logrus.WithError(resperr).Error("Oops")
			return "", resperr
		}

		var boards Boards
		berr := json.Unmarshal([]byte(resp.Body()), &boards)
		if berr != nil {
			fmt.Printf("Error parsing JSON file: %s\n", berr)
		}

		returnValues = append(returnValues, boards.Values...)

		startAt += 50
		if boards.IsLast {
			break
		}
	}

	json, jerr := json.Marshal(returnValues)
	if jerr != nil {
		fmt.Printf("Error parsing JSON file: %s\n", jerr)
	}

	return string(json[:]), nil
}

func (r *Jira) GetActiveSprint(board string) (string, error) {

	// https://alteryx.atlassian.net/rest/agile/1.0/board/388/sprint?state=active
	fetchUri := fmt.Sprintf("%s%s/board/%s/sprint?state=active", r.BaseUrl, r.AgilePath, board)
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		Get(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

func (r *Jira) GetSprintIssues(board string, sprint string) (string, error) {

	startAt := 0
	var returnValues []interface{}

	for {
		// https://alteryx.atlassian.net/rest/agile/1.0/board/388/sprint/723/issue
		fetchUri := fmt.Sprintf("%s%s/board/%s/sprint/%s/issue?startAt=%d", r.BaseUrl, r.AgilePath, board, sprint, startAt)

		resp, resperr := r.Client.R().
			SetHeader("Content-Type", "application/json").
			Get(fetchUri)

		if resperr != nil {
			logrus.WithError(resperr).Error("Oops")
			return "", resperr
		}
		logrus.Warn(resp.Body()[:])
		var sprintIssues SprintIssues
		berr := json.Unmarshal([]byte(resp.Body()), &sprintIssues)
		if berr != nil {
			fmt.Printf("Error parsing JSON file: %s\n", berr)
		}

		returnValues = append(returnValues, sprintIssues.Issues)

		startAt += 50
		if sprintIssues.Total < sprintIssues.MaxResults {
			break
		}
		// if sprintIssues.IsLast {
		// 	break
		// }
	}

	json, jerr := json.Marshal(returnValues)
	if jerr != nil {
		fmt.Printf("Error parsing JSON file: %s\n", jerr)
	}

	return string(json[:]), nil
}
