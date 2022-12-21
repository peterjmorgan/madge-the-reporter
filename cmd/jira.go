package cmd

import (
	"fmt"
	phylum "github.com/peterjmorgan/go-phylum"
	"github.com/peterjmorgan/madge-the-reporter/jira"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(jiraCmd)
	jiraCmd.Flags().StringP("PhylumProjectID", "p", "", "Reference the Phylum Project with ID")
	jiraCmd.Flags().StringP("JiraProjectKey", "j", "", "Reference the Jira Project Key")
}

var myOpts jira.JiraClientOpts = jira.JiraClientOpts{
	OnPrem:      true,
	AuthType:    "PAT",
	Domain:      "http://vader.lan:8080",
	Username:    "pmorgan",
	Token:       os.Getenv("JIRA_PAT"),
	ProjectName: "Vulnerabilities",
	VulnType:    "Vulnerability",
}

var pOpts phylum.ClientOptions = phylum.ClientOptions{Token: os.Getenv("PHYLUM_TOKEN")}

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "jira",
	Run: func(cmd *cobra.Command, args []string) {
		phylumProjectID, _ := cmd.Flags().GetString("PhylumProjectID")
		jiraProjectKey, _ := cmd.Flags().GetString("JiraProjectKey")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		_ = phylumProjectID

		// create clients
		j, err := jira.NewJiraOnPremClient(myOpts)
		if err != nil {
			log.Errorf("failed to create jira client: %v\n", err)
			return
		}
		p, err := phylum.NewClient(&pOpts)
		if err != nil {
			log.Errorf("failed to create phylum client: %v\n", err)
			return
		}

		jiraIssues, err := j.GetJiraIssuesByProject(jiraProjectKey)
		if err != nil {
			log.Errorf("failed to get jira issues: %v\n", err)
			return
		}

		_ = jiraIssues
		//TODO: debug, remove
		//for _, elem := range jiraIssues {
		//	fmt.Printf("%v\n", elem.Fields.Summary)
		//}

		phylumProjectIssues, err := p.GetProjectIssues(phylumProjectID)
		if err != nil {
			log.Errorf("failed to get phylum issues: %v\n", err)
			return
		}

		_ = phylumProjectIssues
		////TODO: debug, remove
		//for _, elem := range phylumProjectIssues {
		//	fmt.Printf("%v\n", elem.Title)
		//}

		// Assumption: Phylum Title -> Summary field
		var ToBeAdded []phylum.IssuesListItem
		for _, pi := range phylumProjectIssues {
			title := pi.Title
			matched := false
			for _, ji := range jiraIssues {
				if title == ji.Fields.Summary {
					matched = true
				}
			}
			if !matched {
				ToBeAdded = append(ToBeAdded, pi)
			}
		}

		if dryRun {
			fmt.Printf("Issues to be added: \n")
			for _, elem := range ToBeAdded {
				fmt.Printf("%v\n", elem.Title)
			}
		} else {
			for _, elem := range ToBeAdded {
				issueKey, err := j.CreateIssue(elem, jiraProjectKey)
				if err != nil {
					log.Errorf("failed to create issue: %v\n", elem.Title)
				} else {
					fmt.Printf("Created %s for %v\n", issueKey, elem.Title)
				}
			}
		}
	},
}
