package jira

import (
	jiraonprem "github.com/andygrunwald/go-jira/v2/onpremise"
	"reflect"
	"testing"
)

var myOpts JiraClientOpts = JiraClientOpts{
	OnPrem:      true,
	AuthType:    "PAT",
	Domain:      "http://vader.lan:8080",
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
		want    []jiraonprem.Issue
		wantErr bool
	}{
		{"one", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := j.getJiraIssuesByProject()
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
