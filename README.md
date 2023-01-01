# Madge

Madge is a reporter of issues, Phylum supply chain security issues specifically. Madge translates Phylum issues into tickets for various ticketing systems. Currently, Madge supports:
- Jira

(other ticketing systems coming soon)

It is a binary with subcommands:
- `configure` - interactively create a configuration file `madge_config.yaml`
- `jira` - parse `madge_config.yaml` to create Jira Issues from Phylum Issues

## Use case
Alice astutely understands using Phylum instead of an SCA tool enables her to defend her organization's software product from the totality of software supply chain risk. Alice wants to report the Phylum Issues for her project to Bob, her developer co-worker, so he can work on fixing the issues. The only issue is, Bob uses Jira for all software development processes. Madge can help here!

1. Alice `analyzes` her software project with Phylum.
1. Alice triages the resulting Phylum Issues, resulting in some *suppressed issues* in her Phylum project.
1. Alice wants to report the remaining Phylum Issues into Jira, so her developer team member, Bob, can fix the issues.
1. `madge` enables Alice to do this!

## Installation
If have you a modern Go toolchain installed:

`go install github.com/peterjmorgan/madge-the-reporter@latest`

## Getting Started
Interactively configure madge:

`./madge configure`

Synchronize Phylum Issues in PHYLUM_PROJECT_ID with a Jira project identified by JIRA_PROJECT_KEY

`./madge jira -j JIRA_PROJECT_KEY -p PHYLUM_PROJECT_ID`

## Configuration
`madge` uses a YAML config file named `madge_config.yaml`. This file can be edited directly, but it is suggested the first configuration be created using the `./madge configure` subcommand.

`madge` looks for the config file in the current working directory. This can be overridden using the `-c PATH_TO_CONFIG` flag. 

### Phylum -> Jira
By default, madge creates Jira issues using the 'Bug' Issue Type in Jira. This can be defined in the madge configuration.

| Phylum Field | Jira Field | Configurable? |
|-----|-----| ----- | 
| Issue Title | Summary | No |
| Issue Description | Description | No |
| Issue Impact (severity) | Undefined | Yes |
| Vulnerability CWE | Undefined | Yes |
| Recommendation/Remediation | Undefined | Yes |

## Workflow
1. User checks out a software project locally (git et al.)
1. User analyzes a software project with Phylum (phylum analyze, or CI automation)
1. User triages Phylum Issues in Phylum Web UI
1. User invokes `madge configure` in the local software project directory to create `madge_config.yaml`
1. User dry-runs madge using the `-D` flag to ensure the everything is properly configured, and the issues to be added look correct
1. User runs madge to add Phylum Issues to the ticketing system
1. User commits `madge_config.yaml` to the VCS system (git et al.) so other users can repeat the process when necessary
1. Developers use the ticketing system to fix issues
1. Phylum analyze is ran again on the updated software project iteratively to confirm issues are fixed
1. Madge is ran again to repopulate the ticketing system with remaining and/or new Phylum Issues




