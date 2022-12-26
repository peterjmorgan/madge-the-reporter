package cmd

import (
	phylum "github.com/peterjmorgan/go-phylum"
	"github.com/peterjmorgan/madge-the-reporter/jira"
	"github.com/peterjmorgan/madge-the-reporter/structs"
	"github.com/peterjmorgan/madge-the-reporter/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(jiraCmd)
	jiraCmd.Flags().StringP("PhylumProjectID", "p", "", "Reference the Phylum Project with ID")
	jiraCmd.Flags().StringP("JiraProjectKey", "j", "", "Reference the Jira Project Key")
}

// var myOpts jira.JiraClientOpts = jira.JiraClientOpts{
// 	OnPrem:      true,
// 	AuthType:    "PAT",
// 	Domain:      "http://vader.lan:8080",
// 	Username:    "pmorgan",
// 	Token:       os.Getenv("JIRA_PAT"),
// 	ProjectName: "Vulnerabilities",
// 	VulnType:    "Vulnerability",
// }

//var pOpts phylum.ClientOptions = phylum.ClientOptions{Token: os.Getenv("PHYLUM_TOKEN")}

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "jira",
	Run: func(cmd *cobra.Command, args []string) {
		phylumProjectID, _ := cmd.Flags().GetString("PhylumProjectID")
		jiraProjectKey, _ := cmd.Flags().GetString("JiraProjectKey")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		debugFlag, _ := cmd.Flags().GetBool("debug")

		if debugFlag {
			log.SetLevel(log.DebugLevel)
		}

		jiraConfig, err := utils.ReadConfigFile(&structs.ConfigFile{})
		if err != nil {
			log.Errorf("failed to read config file: %v\n", err)
			return
		}
		log.Debugf("Read config: %v\n", jiraConfig)

		//TODO: configure should write this structure
		jiraOpts := jira.JiraClientOpts{
			OnPrem:   jiraConfig.JiraConfigObj.OnPrem,
			AuthType: jiraConfig.JiraConfigObj.AuthType,
			Domain:   jiraConfig.JiraConfigObj.URI,
			Username: "blah", // TODO: remove
			Token:    jiraConfig.JiraConfigObj.Token,
			VulnType: "Vulnerability", // TODO: figure this out
			Config:   jiraConfig.JiraConfigObj,
		}

		// create clients
		j, err := jira.NewJiraOnPremClient(jiraOpts)
		if err != nil {
			log.Errorf("failed to create jira client: %v\n", err)
			return
		}
		p, err := phylum.NewClient(&phylum.ClientOptions{jiraConfig.PhylumToken, nil, nil})
		if err != nil {
			log.Errorf("failed to create phylum client: %v\n", err)
			return
		}

		jiraIssues, err := j.GetJiraIssuesByProject(jiraProjectKey)
		if err != nil {
			log.Errorf("failed to get jira issues: %v\n", err)
			return
		}
		log.Infof("Count of Jira Issues in %v: %v\n", jiraProjectKey, len(jiraIssues))

		phylumProjectIssues, err := p.GetProjectIssues(phylumProjectID)
		if err != nil {
			log.Errorf("failed to get phylum issues: %v\n", err)
			return
		}
		log.Infof("Count of Phylum Issues in %v: %v\n", phylumProjectID, len(phylumProjectIssues))

		// Assumption: Phylum Title -> Summary field
		var ToBeAdded []phylum.IssuesListItem
		var ToBeSkipped []phylum.IssuesListItem

		//TODO: matching on PhylumIssue.Title is not a great strategy, should be storing an ID in the metadata
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
			} else {
				ToBeSkipped = append(ToBeSkipped, pi)
			}
		}

		if dryRun {
			log.Infof("DRY-RUN")
			if len(ToBeAdded) > 0 {
				log.Infof("%d issues to be added:\n", len(ToBeAdded))
				for _, elem := range ToBeAdded {
					log.Infof("%v\n", elem.Title)
				}
			}
		} else {
			for _, elem := range ToBeAdded {
				issueKey, err := j.CreateIssue(elem, jiraProjectKey)
				if err != nil {
					log.Errorf("failed to create issue: %v\n", elem.Title)
				} else {
					log.Infof("Created %s for %v\n", issueKey, elem.Title)
				}
			}
		}
		if len(ToBeSkipped) > 0 {
			log.Infof("%d skipped issues (already in Jira Project %v):\n", len(ToBeSkipped), jiraProjectKey)
			for _, elem := range ToBeSkipped {
				log.Infof("%v\n", elem.Title)
			}
		}
	},
}
