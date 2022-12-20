package jira

import (
	"context"
	"fmt"
	jiracloud "github.com/andygrunwald/go-jira/v2/cloud"
	jiraonprem "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/peterjmorgan/go-phylum"
	log "github.com/sirupsen/logrus"
	"github.com/trivago/tgo/tcontainer"
	"net/http"
)

type JiraClientOpts struct {
	OnPrem      bool
	AuthType    string
	Domain      string
	Username    string
	Token       string
	ProjectID   string //TODO: check this
	ProjectName string // needed for queries i think
	VulnType    string
}

type JiraFields struct {
	SeverityField    string
	CweField         string
	DescriptionField string
	ReporterField    string
	SummaryField     string
}

// TODO
// Make a wrapper around NewClient to call NewJiraOnPremClient or NewJiraCloudClient
//func NewJiraClient(opts JiraClientOpts)

type JiraClient struct {
	Client  *jiraonprem.Client
	Opts    JiraClientOpts
	Project jiraonprem.Project
}

func NewJiraOnPremClient(opts JiraClientOpts) (*JiraClient, error) {

	var auth *http.Client
	var jiraClient *jiraonprem.Client

	// onprem
	if opts.AuthType == "PAT" {
		// PAT Token auth
		tempauth := jiraonprem.BearerAuthTransport{
			Token: opts.Token,
		}
		auth = tempauth.Client()

	}

	jiraClient, err := jiraonprem.NewClient(opts.Domain, auth)
	if err != nil {
		log.Errorf("failed to create jira client: %v\n", err)
		return nil, err
	}

	user, _, err := jiraClient.User.GetSelf(context.Background())
	if err != nil {
		log.Errorf("failed to get current user: %v\n", err)
		return nil, err
	}
	_ = user //TODO: debug, remove

	return &JiraClient{
		Client: jiraClient,
		Opts:   opts,
	}, nil
}

// TODO: need to figure out how to make this an interface or something
func NewJiraCloudClient(opts JiraClientOpts) (*jiracloud.Client, error) {
	// cloud
	auth := jiracloud.BasicAuthTransport{
		Username: opts.Username,
		APIToken: opts.Token,
	}

	jiraClient, err := jiracloud.NewClient(opts.Domain, auth.Client())
	if err != nil {
		log.Errorf("failed to create jira client: %v\n", err)
		return nil, err
	}

	user, _, err := jiraClient.User.GetCurrentUser(context.Background())
	if err != nil {
		log.Errorf("failed to get current user: %v\n", err)
		return nil, err
	}
	_ = user //TODO: debug, remove

	return jiraClient, nil
}

// TODO: think about this, most likely will need to pass the project into the query as clients will have multiple projects in Jira
func (j *JiraClient) GetJiraIssuesByProject(projectKey string) ([]jiraonprem.Issue, error) {
	jql := fmt.Sprintf("project = %s and type = %s", projectKey, j.Opts.VulnType)

	issues, _, err := j.Client.Issue.Search(context.Background(), jql, nil)
	if err != nil {
		log.Errorf("failed to search jira for issues: %v\n", err)
		return nil, err
	}

	return issues, nil
}

func (j *JiraClient) CreateIssue(issue phylum.IssuesListItem, projectKey string) (string, error) {

	jiraProject, _, err := j.Client.Project.Get(context.Background(), projectKey)
	if err != nil {
		log.Errorf("failed to get jira project: %v\n", err)
		return "", err
	}

	//unknown := tcontainer.NewMarshalMap()
	severity := tcontainer.NewMarshalMap()
	severity.Set("customfield_10112", "high")

	newIssue := jiraonprem.Issue{
		Fields: &jiraonprem.IssueFields{
			Expand:                        "",
			Type:                          jiraonprem.IssueType{Name: "Vulnerability"},
			Project:                       *jiraProject,
			Environment:                   "",
			Resolution:                    nil,
			Priority:                      nil,
			Resolutiondate:                jiraonprem.Time{},
			Created:                       jiraonprem.Time{},
			Duedate:                       jiraonprem.Date{},
			Watches:                       nil,
			Assignee:                      nil,
			Updated:                       jiraonprem.Time{},
			Description:                   issue.Description,
			Summary:                       issue.Title,
			Creator:                       nil,
			Reporter:                      nil,
			Components:                    nil,
			Status:                        nil,
			Progress:                      nil,
			AggregateProgress:             nil,
			TimeTracking:                  nil,
			TimeSpent:                     0,
			TimeEstimate:                  0,
			TimeOriginalEstimate:          0,
			Worklog:                       nil,
			IssueLinks:                    nil,
			Comments:                      nil,
			FixVersions:                   nil,
			AffectsVersions:               nil,
			Labels:                        nil,
			Subtasks:                      nil,
			Attachments:                   nil,
			Epic:                          nil,
			Sprint:                        nil,
			Parent:                        nil,
			AggregateTimeOriginalEstimate: 0,
			AggregateTimeSpent:            0,
			AggregateTimeEstimate:         0,
			Unknowns:                      nil,
		},
	}

	issueId, resp, err := j.Client.Issue.Create(context.Background(), &newIssue)
	if err != nil {
		log.Errorf("failed to create jira issue: %v\n", err)
		return "", err
	}

	_ = resp
	return issueId.ID, nil
}
