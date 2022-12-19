package structs

type PhylumProject struct {
	Name      string `json:"name" yaml:"name"`
	ID        string `json:"id" yaml:"id"`
	UpdatedAt string `json:"updated_at" yaml:"created_at"`
	Ecosystem string `json:"ecosystem"`
}

// reporting system: jira, github, gitlab, defectdojo?
// token

// jira: domain, email, token
//	project (id)
//  fields

type ConfigThing struct {
	VcsType     string
	VcsToken    string
	Associated  map[string]string
	PhylumToken string
	PhylumGroup string
}

type ConfigFile struct {
	Filename string
}

type JiraConfig struct {
	Domain string
}
