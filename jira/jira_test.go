package jira

import (
	jiraonprem "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/peterjmorgan/go-phylum"
	"reflect"
	"testing"
)

var myOpts JiraClientOpts = JiraClientOpts{
	OnPrem:   true,
	AuthType: "PAT",
	//Domain:      "http://vader.lan:8080",
	Domain:      "http://vader.tail23af1.ts.net:8080",
	Username:    "pmorgan",
	Token:       "NzIwNjIxNzUzMjk5Okn6vIADIjFNaaRxK6NI4o/tU7UP",
	ProjectName: "Vulnerabilities",
	VulnType:    "Vulnerability",
}

func TestNewJiraOnPremClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    JiraClientOpts
		want    *JiraClient
		wantErr bool
	}{
		{"one", myOpts, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJiraOnPremClient(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJiraOnPremClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJiraOnPremClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJiraClient_getJiraIssuesByProject(t *testing.T) {
	j, _ := NewJiraOnPremClient(myOpts)

	tests := []struct {
		name    string
		project string
		want    []jiraonprem.Issue
		wantErr bool
	}{
		{"one", "Vulnerabilities", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := j.GetJiraIssuesByProject(tt.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJiraIssuesByProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getJiraIssuesByProject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJiraClient_CreateIssue(t *testing.T) {
	j, _ := NewJiraOnPremClient(myOpts)
	vid := "VID:000011EE"
	cwe := "IVCWE-601"

	type args struct {
		issue      phylum.IssuesListItem
		projectKey string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"one", args{
			phylum.IssuesListItem{
				Description: "### Overview\nVersions of `serve-static` prior to 1.6.5 ( or 1.7.x prior to 1.7.2 ) are affected by an open redirect vulnerability on some browsers when configured to mount at the root directory.\n\n\n## Proof of Concept\n\nA link to `http://example.com//www.google.com/%2e%2e` will redirect to `//www.google.com/%2e%2e`\n\nSome browsers will interpret this as `http://www.google.com/%2e%2e`, resulting in an external redirect.\n\n\n\n\n### Recommendation\nVersion 1.7.x: Update to version 1.7.2 or later.\nVersion 1.6.x: Update to version 1.6.5 or later.\nUpgrade to one of versions 1.7.2,1.7.2.",
				Id:          &vid,
				Ignored:     phylum.IgnoredReason("false"),
				Impact:      "low",
				RiskType:    "vulnerability",
				Score:       0,
				Tag:         &cwe,
				Title:       "serve-static@1.5.3 is vulnerable to open redirect",
			}, "VULN",
		}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := j.CreateIssue(tt.args.issue, tt.args.projectKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateIssue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
