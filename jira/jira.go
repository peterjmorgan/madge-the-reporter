package jira

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	jiracloud "github.com/andygrunwald/go-jira/v2/cloud"
	jiraonprem "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/peterjmorgan/go-phylum"
	"github.com/peterjmorgan/madge-the-reporter/structs"
	log "github.com/sirupsen/logrus"
	"github.com/trivago/tgo/tcontainer"
	stripmd "github.com/writeas/go-strip-markdown"
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
	Config      structs.JiraConfig
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

//type FieldMapping struct {
//	Recommendation struct {
//		Name string
//		Id string
//	}
//	CWE struct {
//		Name string
//		Id string
//	}
//	Severity struct {
//		Name string
//		Id string
//	}
//}

func (j *JiraClient) CreateIssue(issue phylum.IssuesListItem, projectKey string) (string, error) {

	jiraProject, _, err := j.Client.Project.Get(context.Background(), projectKey)
	if err != nil {
		log.Errorf("failed to get jira project: %v\n", err)
		return "", err
	}

	// Set custom fields
	unknown := tcontainer.NewMarshalMap()
	sev := make(map[string]string, 0)

	// Check if custom severity fields are set
	if j.Opts.Config.CustomFields.Severity.ID != "" && j.Opts.Config.CustomFields.Severity.Name != "" {
		// interpret Phylum impact level as user-configured severity level
		switch strings.ToLower(string(issue.Impact)) {
		case "critical":
			if j.Opts.Config.SeverityFields.Critical.ID != "" {
				// Severity custom field ID is set
				sev["value"] = j.Opts.Config.SeverityFields.Critical.Name
				sev["id"] = j.Opts.Config.SeverityFields.Critical.ID
			}
		case "high":
			if j.Opts.Config.SeverityFields.High.ID != "" {
				// Severity custom field ID is set
				sev["value"] = j.Opts.Config.SeverityFields.High.Name
				sev["id"] = j.Opts.Config.SeverityFields.High.ID
			}
		case "medium":
			if j.Opts.Config.SeverityFields.Medium.ID != "" {
				// Severity custom field ID is set
				sev["value"] = j.Opts.Config.SeverityFields.Medium.Name
				sev["id"] = j.Opts.Config.SeverityFields.Medium.ID
			}
		case "low":
			if j.Opts.Config.SeverityFields.Low.ID != "" {
				// Severity custom field ID is set
				sev["value"] = j.Opts.Config.SeverityFields.Low.Name
				sev["id"] = j.Opts.Config.SeverityFields.Low.ID
			}
		}
	}

	// set the severity field
	// unknown.Set("customfield_10112", sev)
	// set the severity field using the customfield definition, then set the val of the k->v mapping to the sev map
	//unknown.Set(textFieldsMapping["severity"]["id"], sev)
	unknown.Set(j.Opts.Config.CustomFields.Severity.ID, sev)

	// Set CWE
	if issue.RiskType == "vulnerability" {
		cwePat := regexp.MustCompile(`..(CWE-\d\d\d)`)
		if doesMatch := cwePat.MatchString(*issue.Tag); doesMatch {
			matches := cwePat.FindStringSubmatch(*issue.Tag)
			cwe := matches[1]
			if len(cwe) > 0 && j.Opts.Config.CustomFields.CWE.ID != "" {
				//unknown.Set("customfield_10113", cwe)
				//unknown.Set(textFieldsMapping["cwe"]["id"], cwe)
				unknown.Set(j.Opts.Config.CustomFields.CWE.ID, cwe)
			} else {
				log.Errorf("CWE field len = 0 - Phylum Vuln Issue: %v\n", issue.Title)
			}
		}
	}

	// Set recomendation
	if j.Opts.Config.CustomFields.Recommendation.ID != "" {
		recommendation, err := phylum.ExtractRemediation(issue)
		if err != nil {
			log.Errorf("failed to extract remediation from %v: %v\n", issue.Title, err)
		}
		unknown.Set(j.Opts.Config.CustomFields.Recommendation.ID, recommendation)
	}

	// Create issue with fields set
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
