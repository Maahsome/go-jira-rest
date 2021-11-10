package jira

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Jira struct {
	BaseUrl string
	ApiPath string
	User    string
	Token   string
	Client  *resty.Client
}

// New generate a new jira client
func New(baseUrl, apiPath, user, token string) *Jira {

	restClient := resty.New()

	if apiPath == "" {
		apiPath = "/rest/api/3"
	}
	restClient.SetBasicAuth(user, token)

	return &Jira{
		BaseUrl: baseUrl,
		ApiPath: apiPath,
		Token:   token,
		Client:  restClient,
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
