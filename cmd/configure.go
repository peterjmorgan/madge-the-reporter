package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/peterjmorgan/madge-the-reporter/structs"
	"github.com/peterjmorgan/madge-the-reporter/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configure",
	Run: func(cmd *cobra.Command, args []string) {

		configData, err := utils.ReadConfigFile(&structs.ConfigFile{})
		if err == nil {
			configJson, err3 := json.MarshalIndent(configData, "", "  ")
			if err3 != nil {
				fmt.Printf("Failed to unmarshall config data: %v\n", err3)
				return
			}

			fmt.Printf("Found existing config:\n %s\n", string(configJson))
			reconfigurePrompt := promptui.Prompt{
				Label:     "Would you like to enter new configuration data",
				Default:   "n",
				IsConfirm: true,
			}

			shouldReconfigure, err := reconfigurePrompt.Run()
			if err != nil && err.Error() != "" {
				fmt.Printf(err.Error())
				return
			}
			shouldReconfigure = strings.ToLower(shouldReconfigure)
			if shouldReconfigure == "" || shouldReconfigure == "n" {
				fmt.Printf("No configuration change.\n")
				os.Exit(0)
			}
			// shouldReconfigure == "y" will fall through
		} else {
			// doesn't exist
			fmt.Printf("Could not find 'madge_config.yaml', let's create one")
		}

		// Configure issue system type
		issueSystemPrompt := promptui.Select{
			Label: "Choose Issue system",
			Items: []string{"Jira"},
		}

		_, issueSystemResult, err := issueSystemPrompt.Run()
		if err != nil {
			fmt.Printf("issueSystemPrompt failed: %v\n", err)
		}

		//config := make(map[string]string, 0)
		//customFields := make(map[string]map[string]string, 0)
		//severityFields := make(map[string]map[string]string, 0)

		var rootConfig structs.RootConfig

		switch issueSystemResult {
		case "Jira":
			//config["IssueSystem"] = "jira"
			rootConfig.IssueSystem = "jira"

			// Jira Deployment: Onprem or Cloud
			isDeploymentPrompt := promptui.Select{
				Label: "Choose Jira deployment",
				Items: []string{"On Prem"},
			}
			_, isDeploymentResult, err := isDeploymentPrompt.Run()
			if err != nil {
				log.Errorf("isDeploymentPrompt failed: %v\n", err)
				return
			}
			//config["Deployment"] = isDeploymentResult
			if isDeploymentResult == "On Prem" {
				rootConfig.JiraConfigObj.OnPrem = true
			}

			// URL to issue tracker for API use
			isUri, err := utils.PromptForString(fmt.Sprintf("Enter %s API URI including url scheme:", rootConfig.IssueSystem), -1)
			if err != nil {
				log.Errorf("isUrl failed: %v\n", err)
				return
			}
			//config["URI"] = isUri
			rootConfig.JiraConfigObj.URI = isUri

			// Authentication Type:
			isAuthenticationType := promptui.Select{
				Label: fmt.Sprintf("Choose authentication type for %v\n", rootConfig.IssueSystem),
				Items: []string{"PAT"}, // "Oauth",

			}
			_, isAuthenticationTypeResult, err := isAuthenticationType.Run()
			if err != nil {
				log.Errorf("isAuthenticationType failed: %v\n", err)
				return
			}
			//config["AuthenticationType"] = isAuthenticationTypeResult
			rootConfig.JiraConfigObj.AuthType = isAuthenticationTypeResult

			//// Username:
			//isUsername, err := utils.PromptForString(fmt.Sprintf("Enter username for %v API:", config["IssueSystem"]), -1)
			//if err != nil {
			//	log.Errorf("isUsername failed: %v\n", err)
			//	return
			//}
			//config["Username"] = isUsername

			// Jira Token:
			isToken, err := utils.PromptForString(fmt.Sprintf("Enter Personal Access Token for %s API:", rootConfig.IssueSystem), -1)
			if err != nil {
				log.Errorf("isToken failed: %v\n", err)
				return
			}
			//config["JiraToken"] = isToken
			rootConfig.JiraConfigObj.Token = isToken

			// custom issue fields
			customIssueFieldPrompt := promptui.Prompt{
				Label:     "Madge can write Phylum data to custom issue fields in Jira. Available fields (and field type) include: Severity (select chooser), Remediation/Recommendation (text field), CWE (text field), Would you like to configure Madge to use custom fields?",
				Default:   "n",
				IsConfirm: true,
			}

			shouldConfigureCustomIssueFields, err := customIssueFieldPrompt.Run()
			if err != nil && err.Error() != "" {
				log.Errorf(err.Error())
				return
			}
			if strings.ToLower(shouldConfigureCustomIssueFields) == "y" {
				// Severity
				severityCustomIssueFieldPrompt := promptui.Prompt{
					Label:     "Configure Severity? This can be configured only as a select box",
					Default:   "n",
					IsConfirm: true,
				}

				shouldConfigureSeverity, err := severityCustomIssueFieldPrompt.Run()
				if err != nil && err.Error() != "" {
					log.Errorf(err.Error())
					return
				}
				if strings.ToLower(shouldConfigureSeverity) == "y" {

					severityCustomFieldName, err := utils.PromptForString("Enter the Severity chooser field NAME", -1)
					if err != nil {
						log.Errorf("severityCustomFieldName failed: %v\n", err)
						return
					}
					severityCustomFieldId, err := utils.PromptForString("Enter the Severity field ID", -1)
					if err != nil {
						log.Errorf("severityCustomFieldName failed: %v\n", err)
						return
					}

					//severityFields["severity"] = make(map[string]string,0)
					//severityFields["severity"]["value"] = severityCustomFieldName
					//severityFields["severity"]["id"] = severityCustomFieldId
					rootConfig.JiraConfigObj.CustomFields.Severity.Name = severityCustomFieldName
					rootConfig.JiraConfigObj.CustomFields.Severity.ID = severityCustomFieldId

					severityCriticalCustomFieldName, err := utils.PromptForString("Enter the Critical severity field NAME", -1)
					if err != nil {
						log.Errorf("severityCriticalCustomFieldName failed: %v\n", err)
						return
					}
					severityCriticalCustomFieldId, err := utils.PromptForString("Enter the Critical severity field ID", -1)
					if err != nil {
						log.Errorf("severityCustomFieldName failed: %v\n", err)
						return
					}
					//severityFields["critical"] = make(map[string]string,0)
					//severityFields["critical"]["name"] = severityCriticalCustomFieldName
					//severityFields["critical"]["id"] = severityCriticalCustomFieldId
					rootConfig.JiraConfigObj.SeverityFields.Critical.Name = severityCriticalCustomFieldName
					rootConfig.JiraConfigObj.SeverityFields.Critical.ID = severityCriticalCustomFieldId

					// high
					severityHighCustomFieldName, err := utils.PromptForString("Enter the High severity field NAME", -1)
					if err != nil {
						log.Errorf("severityHighCustomFieldName failed: %v\n", err)
						return
					}
					severityHighCustomFieldId, err := utils.PromptForString("Enter the High severity field ID", -1)
					if err != nil {
						log.Errorf("severityHighCustomFieldName failed: %v\n", err)
						return
					}
					//severityFields["high"] = make(map[string]string,0)
					//severityFields["high"]["name"] = severityHighCustomFieldName
					//severityFields["high"]["id"] = severityHighCustomFieldId
					rootConfig.JiraConfigObj.SeverityFields.High.Name = severityHighCustomFieldName
					rootConfig.JiraConfigObj.SeverityFields.High.ID = severityHighCustomFieldId

					// Medium
					severityMediumCustomFieldName, err := utils.PromptForString("Enter the Medium severity field NAME", -1)
					if err != nil {
						log.Errorf("severityMediumCustomFieldName failed: %v\n", err)
						return
					}
					severityMediumCustomFieldId, err := utils.PromptForString("Enter the Medium severity field ID", -1)
					if err != nil {
						log.Errorf("severityMediumCustomFieldName failed: %v\n", err)
						return
					}
					//severityFields["medium"] = make(map[string]string,0)
					//severityFields["medium"]["name"] = severityMediumCustomFieldName
					//severityFields["medium"]["id"] = severityMediumCustomFieldId
					rootConfig.JiraConfigObj.SeverityFields.Medium.Name = severityMediumCustomFieldName
					rootConfig.JiraConfigObj.SeverityFields.Medium.ID = severityMediumCustomFieldId

					// Low
					severityLowCustomFieldName, err := utils.PromptForString("Enter the Low severity field NAME", -1)
					if err != nil {
						log.Errorf("severityLowCustomFieldName failed: %v\n", err)
						return
					}
					severityLowCustomFieldId, err := utils.PromptForString("Enter the Low severity field ID", -1)
					if err != nil {
						log.Errorf("severityLowCustomFieldName failed: %v\n", err)
						return
					}
					//severityFields["low"] = make(map[string]string,0)
					//severityFields["low"]["name"] = severityLowCustomFieldName
					//severityFields["low"]["id"] = severityLowCustomFieldId
					rootConfig.JiraConfigObj.SeverityFields.Low.Name = severityLowCustomFieldName
					rootConfig.JiraConfigObj.SeverityFields.Low.ID = severityLowCustomFieldId
				}
				// Remediation/Recommendation
				recommendationCustomIssueFieldPrompt := promptui.Prompt{
					Label:     "Configure Recommendation text field?",
					Default:   "n",
					IsConfirm: true,
				}

				shouldConfigureRecommendation, err := recommendationCustomIssueFieldPrompt.Run()
				if err != nil && err.Error() != "" {
					log.Errorf(err.Error())
					return
				}

				if strings.ToLower(shouldConfigureRecommendation) == "y" {
					recommendationCustomFieldName, err := utils.PromptForString("Enter the Recommendation/Remediation field NAME", -1)
					if err != nil {
						log.Errorf("recommendationCustomFieldName failed: %v\n", err)
						return
					}
					recommendationCustomFieldId, err := utils.PromptForString("Enter the Recommendation/Remediation field ID", -1)
					if err != nil {
						log.Errorf("recommendationCustomFieldId failed: %v\n", err)
						return
					}
					//customFields["recommendation"] = make(map[string]string,0)
					//customFields["recommendation"]["value"] = recommendationCustomFieldName
					//customFields["recommendation"]["id"] = recommendationCustomFieldId
					rootConfig.JiraConfigObj.CustomFields.Recommendation.Name = recommendationCustomFieldName
					rootConfig.JiraConfigObj.CustomFields.Recommendation.ID = recommendationCustomFieldId
				}
				// CWE
				cweCustomIssueFieldPrompt := promptui.Prompt{
					Label:     "Configure CWE text field?",
					Default:   "n",
					IsConfirm: true,
				}

				shouldConfigureCWE, err := cweCustomIssueFieldPrompt.Run()
				if err != nil && err.Error() != "" {
					log.Errorf(err.Error())
					return
				}

				if strings.ToLower(shouldConfigureCWE) == "y" {
					cweCustomFieldName, err := utils.PromptForString("Enter the CWE field NAME", -1)
					if err != nil {
						log.Errorf("cweCustomFieldName failed: %v\n", err)
						return
					}
					cweCustomFieldId, err := utils.PromptForString("Enter the CWE field ID", -1)
					if err != nil {
						log.Errorf("cweCustomFieldId failed: %v\n", err)
						return
					}
					//customFields["cwe"] = make(map[string]string,0)
					//customFields["cwe"]["value"] = cweCustomFieldName
					//customFields["cwe"]["id"] = cweCustomFieldId
					rootConfig.JiraConfigObj.CustomFields.CWE.Name = cweCustomFieldName
					rootConfig.JiraConfigObj.CustomFields.CWE.ID = cweCustomFieldId
				}
			}

		default:
			fmt.Printf("issueSystemType switch default case. Shouldn't happen\n")
			return
		}

		// Configure Phylum token
		phylumConfigurePrompt := promptui.Prompt{
			Label:     "Madge can try to detect your Phylum token. Would you like to OVERRIDE and set it manually",
			Default:   "n",
			IsConfirm: true,
		}

		manuallyConfigPhylum, err := phylumConfigurePrompt.Run()
		if err != nil && err.Error() != "" {
			fmt.Printf(err.Error())
			return
		}
		var phylumToken string
		if strings.ToLower(manuallyConfigPhylum) == "y" {
			phylumToken, err = utils.PromptForString("Enter Phylum Token (result of `phylum auth token`)", -1)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
		} else {
			phylumToken, err = utils.PhylumGetAuthToken()
			if err != nil {
				fmt.Printf("Failed to read phylum auth token from locally-installed 'phylum'")
				return
			}
			fmt.Printf("Found phylum token from locally-installed 'phylum':\n%v\n", phylumToken)
		}
		//config["PHYLUM_API_KEY"] = phylumToken
		rootConfig.PhylumToken = phylumToken

		// Configure Phylum API URL
		phylumConfigureUrlPrompt := promptui.Prompt{
			Label:     "Are you using a self-hosted instance of Phylum?",
			Default:   "n",
			IsConfirm: true,
		}

		selfHostedPhylum, err := phylumConfigureUrlPrompt.Run()
		if err != nil && err.Error() != "" {
			fmt.Printf(err.Error())
			return
		}
		if strings.ToLower(selfHostedPhylum) == "y" {
			// URL to Phylum API
			isPhylumUri, err := utils.PromptForString("Enter Phylum API URI including url scheme:", -1)
			if err != nil {
				log.Errorf("isPhylumUrl failed: %v\n", err)
				return
			}
			//config["URI"] = isUri
			rootConfig.PhylumHost = isPhylumUri
		}

		yamlData, err := yaml.Marshal(rootConfig)
		if err != nil {
			fmt.Printf("Failed to marshall config data: %v\n", err)
			return
		}

		err = os.WriteFile("madge_config.yaml", yamlData, 0644)
		if err != nil {
			fmt.Printf("Failed to write config file: %v\n", err)
			return
		}

		fmt.Printf("Finished configuring Madge. Wrote configuration to 'madge_config.yaml`\n")
	},
}
