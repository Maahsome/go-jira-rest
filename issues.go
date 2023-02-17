package jira

type Issues struct {
	Expand     string        `json:"expand"`
	Issues     []IssuesIssue `json:"issues"`
	MaxResults int           `json:"maxResults"`
	Names      IssuesNames   `json:"names"`
	Schema     IssuesSchema  `json:"schema"`
	StartAt    int           `json:"startAt"`
	Total      int           `json:"total"`
	IsLast     bool          `json:"isLast"`
}

type IssuesIssue struct {
	Expand     string                `json:"expand"`
	Fields     IssuesIssueFields     `json:"fields"`
	ID         string                `json:"id"`
	Key        string                `json:"key"`
	Operations IssuesIssueOperations `json:"operations"`
	Self       string                `json:"self"`
}

type IssuesIssueFields struct {
	Assignee IssuesIssueFieldsAssignee `json:"assignee"`
	Status   IssuesIssueFieldsStatus   `json:"status"`
	Summary  string                    `json:"summary"`
}

type IssuesIssueFieldsAssignee struct {
	AccountID    string                              `json:"accountId"`
	AccountType  string                              `json:"accountType"`
	Active       bool                                `json:"active"`
	AvatarUrls   IssuesIssueFieldsAssigneeAvatarUrls `json:"avatarUrls"`
	DisplayName  string                              `json:"displayName"`
	EmailAddress string                              `json:"emailAddress"`
	Self         string                              `json:"self"`
	TimeZone     string                              `json:"timeZone"`
}

type IssuesIssueFieldsAssigneeAvatarUrls struct {
	X16x16 string `json:"16x16"`
	X24x24 string `json:"24x24"`
	X32x32 string `json:"32x32"`
	X48x48 string `json:"48x48"`
}

type IssuesIssueFieldsStatus struct {
	Description    string                                `json:"description"`
	IconURL        string                                `json:"iconUrl"`
	ID             string                                `json:"id"`
	Name           string                                `json:"name"`
	Self           string                                `json:"self"`
	StatusCategory IssuesIssueFieldsStatusStatusCategory `json:"statusCategory"`
}

type IssuesIssueFieldsStatusStatusCategory struct {
	ColorName string `json:"colorName"`
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Self      string `json:"self"`
}

type IssuesIssueOperations struct {
	LinkGroups []IssuesIssueOperationsLinkGroup `json:"linkGroups"`
}

type IssuesIssueOperationsLinkGroup struct {
	Groups []IssuesIssueOperationsLinkGroupGroup `json:"groups"`
	ID     string                                `json:"id"`
	Links  []interface{}                         `json:"links"`
}

type IssuesIssueOperationsLinkGroupGroup struct {
	Groups []interface{}                             `json:"groups"`
	ID     string                                    `json:"id"`
	Links  []IssuesIssueOperationsLinkGroupGroupLink `json:"links"`
	Weight int                                       `json:"weight"`
}

type IssuesIssueOperationsLinkGroupGroupLink struct {
	Href       string `json:"href"`
	IconClass  string `json:"iconClass"`
	ID         string `json:"id"`
	Label      string `json:"label"`
	StyleClass string `json:"styleClass"`
	Title      string `json:"title"`
	Weight     int    `json:"weight"`
}

type IssuesNames struct {
	Assignee string `json:"assignee"`
	Status   string `json:"status"`
	Summary  string `json:"summary"`
}

type IssuesSchema struct {
	Assignee IssuesSchemaAssignee `json:"assignee"`
	Status   IssuesSchemaStatus   `json:"status"`
	Summary  IssuesSchemaSummary  `json:"summary"`
}

type IssuesSchemaAssignee struct {
	System string `json:"system"`
	Type   string `json:"type"`
}

type IssuesSchemaStatus struct {
	System string `json:"system"`
	Type   string `json:"type"`
}

type IssuesSchemaSummary struct {
	System string `json:"system"`
	Type   string `json:"type"`
}
