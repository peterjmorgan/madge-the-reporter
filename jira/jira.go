package jira

import (
	"context"
	"fmt"
	jiracloud "github.com/andygrunwald/go-jira/v2/cloud"
	jiraonprem "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/peterjmorgan/go-phylum"
	log "github.com/sirupsen/logrus"
	"github.com/trivago/tgo/tcontainer"
	stripmd "github.com/writeas/go-strip-markdown"
	"io"
	"net/http"
	"regexp"
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

//func ConvertMarkdown(description string) (string, error) {
//	var result string = ""
//	pdoc, err := pandoc.New(nil)
//	if err != nil {
//		log.Errorf("failed to create pandoc client: %v\n", err)
//	}
//	pdoc.
//	return result, nil
//}

func (j *JiraClient) CreateIssue(issue phylum.IssuesListItem, projectKey string) (string, error) {

	jiraProject, _, err := j.Client.Project.Get(context.Background(), projectKey)
	if err != nil {
		log.Errorf("failed to get jira project: %v\n", err)
		return "", err
	}

	// Convert Phylum Issue Description from Markdown to Jira Textile

	// Set custom fields
	unknown := tcontainer.NewMarshalMap()
	sev := make(map[string]interface{}, 0)
	sev["value"] = "critical"
	sev["id"] = "10003"
	unknown.Set("customfield_10112", sev)

	// Set CWE
	if issue.RiskType == "vulnerability" {
		cwePat := regexp.MustCompile(`..(CWE-\d\d\d)`)
		if doesMatch := cwePat.MatchString(*issue.Tag); doesMatch {
			matches := cwePat.FindStringSubmatch(*issue.Tag)
			cwe := matches[1]
			unknown.Set("customfield_10113", cwe)
		}
	}

	newIssue := jiraonprem.Issue{
		Fields: &jiraonprem.IssueFields{
			Expand:         "",
			Type:           jiraonprem.IssueType{Name: "Vulnerability"},
			Project:        *jiraProject,
			Resolutiondate: jiraonprem.Time{},
			Created:        jiraonprem.Time{},
			Duedate:        jiraonprem.Date{},
			Watches:        nil,
			Assignee:       nil,
			Updated:        jiraonprem.Time{},
			Description:    stripmd.Strip(issue.Description),
			Summary:        issue.Title,
			Creator:        nil,
			Reporter:       nil,
			Unknowns:       unknown,
		},
	}

	issueId, resp, err := j.Client.Issue.Create(context.Background(), &newIssue)
	if err != nil {
		log.Errorf("failed to create jira issue: %v\n", err)
		tempbody, _ := io.ReadAll(resp.Body)
		_ = tempbody
		return "", err
	}

	_ = resp
	return issueId.ID, nil
}
