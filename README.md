# Madge

Madge is a reporter of issues, Phylum supply chain security issues specifically. Madge translates Phylum issues into tickets for various ticketing systems. Currently, Madge supports:
- Jira

(other ticketing systems coming soon)

It is a binary with subcommands:
- `configure` - interactively create a configuration file `madge_config.yaml`
- `jira` - parse `madge_config.yaml` to create Jira Issues from Phylum Issues

## Workflow
Alice astutely understands using Phylum instead of an SCA tool enables her to defend her organization's software product from the totality of software supply chain risk. Alice wants to report the Phylum Issues for her project to Bob, her developer co-worker, so he can work on fixing the issues. The only issue is, Bob uses Jira for all software development processes. Madge can help here!

1. Alice `analyzes` her software project with Phylum.
1. Alice triages the resulting Phylum Issues, resulting in some *suppressed issues* in her Phylum project.
1. Alice wants to report the remaining Phylum Issues into Jira, so her developer team member, Bob, can fix the issues.
1. `madge` enables Alice to do this!

## Installation
If have you a modern Go toolchain installed:
`go install github.com/peterjmorgan/madge-the-reporter@latest`

## Getting Started
Synchronize Phylum Issues in PHYLUM_PROJECT_ID with a Jira project identified by JIRA_PROJECT_KEY

`./madge jira -j JIRA_PROJECT_KEY -p PHYLUM_PROJECT_ID`

Interactively configure madge:

`./madge configure`

## Configuration
`madge` uses a YAML configuration file named `madge_config.yaml`. This file can be edited directly, but it is suggested the first configuration be created using the `./madge configure` subcommand.

### Phylum -> Jira
By default, madge creates Jira issues using the 'Bug' Issue Type in Jira. This can be defined in the madge configuration.

| Phylum Field | Jira Field | Configurable? |
|-----|-----| ----- | 
| Issue Title | Summary | No |
| Issue Description | Description | No |
| Issue Impact (severity) | Undefined | Yes |
| Vulnerability CWE | Undefined | Yes |
| Recommendation/Remediation | Undefined | Yes |





