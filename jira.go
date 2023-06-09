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
	Expand     string        `json:"expand"`
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
	IsLast     bool          `json:"isLast"`
	Total      int           `json:"total"`
	Issues     []interface{} `json:"issues"`
}

type Comment struct {
	Body CommentBody `json:"body"`
}

type CommentBody struct {
	Content []CommentBodyContent `json:"content"`
	Type    string               `json:"type"`
	Version int                  `json:"version"`
}

type CommentBodyContent struct {
	Content []CommentBodyContentContent `json:"content"`
	Type    string                      `json:"type"`
}

type CommentBodyContentContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
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
	commentTemplate := `{
	"body": {
		"type": "doc",
		"version": 1,
		"content": [
		  {
			"type": "paragraph",
			"content": [
			  {
				"text": "%s",
				"type": "text"
			  }
			]
		  }
		]
	  }
	}`
	body := fmt.Sprintf(commentTemplate, comment)
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

func (r *Jira) AddCommentMulti(issue string, comment *Comment) (string, error) {

	fetchUri := fmt.Sprintf("%s%s/issue/%s/comment", r.BaseUrl, r.ApiPath, issue)

	body, merr := json.Marshal(comment)
	if merr != nil {
		return "", merr
	}
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

func (r *Jira) SetDescription(issue string, description string) (bool, error) {

	fetchUri := fmt.Sprintf("%s%s/issue/%s", r.BaseUrl, r.ApiPath, issue)
	descriptionTemplate := `{
	"fields": {
	  "description": {
		  "type": "doc",
		  "version": 1,
		  "content": [
		    {
			  "type": "paragraph",
			  "content": [
			    {
				  "text": "%s",
				  "type": "text"
			    }
			  ]
		    }
		  ]
	    }
	  }
	}`

	body := fmt.Sprintf(descriptionTemplate, description)
	resp, resperr := r.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(fetchUri)

	if resperr != nil {
		logrus.WithError(resperr).Error("Oops")
		return false, resperr
	}

	return resp.IsSuccess(), nil
}

// func (r *Jira) SetDescriptionMulti(issue string, comment *Comment) (string, error) {

// 	fetchUri := fmt.Sprintf("%s%s/issue/%s", r.BaseUrl, r.ApiPath, issue)

// 	body, merr := json.Marshal(comment)
// 	if merr != nil {
// 		return "", merr
// 	}
// 	resp, resperr := r.Client.R().
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(body).
// 		Put(fetchUri)

// 	if resperr != nil {
// 		logrus.WithError(resperr).Error("Oops")
// 		return "", resperr
// 	}

// 	return string(resp.Body()[:]), nil
// }

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

func (r *Jira) GetIssuesViaJQL(jql string) ([]IssuesIssue, error) {

	// 	curl --request POST \
	//   --url 'https://alteryx.atlassian.net/rest/api/3/search' \
	//   --user 'christopher.maahs@alteryx.com:${JIRA_TOKEN}' \
	//   --header 'Accept: application/json' \
	//   --header 'Content-Type: application/json' \
	//   --data '{
	//   "expand": [
	//     "names",
	//     "schema",
	//     "operations"
	//   ],
	//   "fields": [
	//     "summary",
	//     "status",
	//     "assignee"
	//   ],
	//   "fieldsByKeys": false,
	//   "jql": "project = TSAASPD",
	//   "maxResults": 15,
	//   "startAt": 0
	// }'

	startAt := 0
	var returnValues []IssuesIssue

	searchTpl := `{
  "expand": [
    "names",
    "schema",
    "operations"
  ],
  "fields": [
    "summary",
    "status",
    "assignee"
  ],
  "fieldsByKeys": false,
  "jql": "%s",
  "maxResults": 15,
  "startAt": %d
}`

	for {
		fetchUri := fmt.Sprintf("%s%s/search", r.BaseUrl, r.ApiPath)

		body := fmt.Sprintf(searchTpl, jql, startAt)
		resp, resperr := r.Client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(fetchUri)

		if resperr != nil {
			logrus.WithError(resperr).Error("Oops")
			return returnValues, resperr
		}

		var issues Issues
		berr := json.Unmarshal([]byte(resp.Body()), &issues)
		if berr != nil {
			// logrus.WithError(berr).Error("Error parsing ISSUES JSON data")
			fmt.Printf("Error parsing ISSUES JSON file: %s\n", berr)
		}

		returnValues = append(returnValues, issues.Issues...)

		startAt += 15
		if issues.Total < issues.MaxResults {
			break
		}
		if issues.IsLast {
			break
		}
	}

	// json, jerr := json.Marshal(returnValues)
	// if jerr != nil {
	// 	fmt.Printf("Error parsing JSON file: %s\n", jerr)
	// }

	return returnValues, nil
}

func (r *Jira) GetSprintIssues(board string, sprint int) (string, error) {

	startAt := 0
	var returnValues []interface{}

	for {
		// https://alteryx.atlassian.net/rest/agile/1.0/board/388/sprint/723/issue
		fetchUri := fmt.Sprintf("%s%s/board/%s/sprint/%d/issue?startAt=%d", r.BaseUrl, r.AgilePath, board, sprint, startAt)

		resp, resperr := r.Client.R().
			SetHeader("Content-Type", "application/json").
			Get(fetchUri)

		if resperr != nil {
			logrus.WithError(resperr).Error("Oops")
			return "", resperr
		}
		var sprintIssues SprintIssues
		berr := json.Unmarshal([]byte(resp.Body()), &sprintIssues)
		if berr != nil {
			fmt.Printf("Error parsing JSON file: %s\n", berr)
		}

		returnValues = append(returnValues, sprintIssues.Issues...)

		startAt += 50
		if sprintIssues.Total < sprintIssues.MaxResults {
			break
		}
		if sprintIssues.IsLast {
			break
		}
	}

	jsonData, jerr := json.Marshal(returnValues)
	if jerr != nil {
		fmt.Printf("Error parsing JSON file: %s\n", jerr)
	}

	issueList := fmt.Sprintf("{ \"issues\": %s }", jsonData)
	return issueList, nil
}
