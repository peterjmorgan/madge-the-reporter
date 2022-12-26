package structs

type PhylumProject struct {
	Name      string `json:"name" yaml:"name"`
	ID        string `json:"id" yaml:"id"`
	UpdatedAt string `json:"updated_at" yaml:"created_at"`
	Ecosystem string `json:"ecosystem"`
}

type ConfigFile struct {
	Filename string
}

type TestConfigData struct {
	Filename string
}

// type ConfigThing struct {
// 	VcsType     string
// 	VcsToken    string
// 	Associated  map[string]string
// 	PhylumToken string
// 	PhylumGroup string
// }
// type MadgeConfig struct {
// 	Config map[string]string
// 	CustomFields map[string]map[string]string
// 	SeverityField map[string]map[string]string
// }

// TODO: Coming up with a better way to do this....
// RootConfig is the struct the config file will unmarshall to, and be marshalled from
type RootConfig struct {
	IssueConfig
	IssueSystem string
	PhylumToken string
	PhylumHost  string
}

// IssueConfig is the holding cell type for JiraConfig and Other specific Configs
// This will be selected and used by what is defined in the RootConfig.IssueSystem
type IssueConfig struct {
	JiraConfigObj       JiraConfig
	DefectDojoConfigObj DefectDojoConfig
}

// main Jira Config structure
type JiraConfig struct {
	OnPrem         bool
	AuthType       string // can only be "PAT" for now
	URI            string
	Token          string
	IssueTypeID    string
	CustomFields   JiraCustomFields
	SeverityFields JiraSeverityConfig
}

// Custom Field structure for JiraConfig
type JiraCustomFields struct {
	CWE struct {
		Name string
		ID   string
	}
	Severity struct {
		Name string
		ID   string
	}
	Recommendation struct {
		Name string
		ID   string
	}
}

// Custom Severity field structure for JiraConfig
// This maps Phylum's impact to the user's severity
// eg. "A Phylum Critical issue maps to the user's 'Catastrophic' severity"
type JiraSeverityConfig struct {
	Critical struct { // Critical is the Phylum impact value
		Name string // Name is the user-provided name of the field for Phylum's "critical"
		ID   string
	}
	High struct {
		Name string
		ID   string
	}
	Medium struct {
		Name string
		ID   string
	}
	Low struct {
		Name string
		ID   string
	}
}

// DefectDojoConfig just to simulate what having another issue system config would look like
type DefectDojoConfig struct {
	OnPrem    bool
	Domain    string
	Username  string
	Passwword string
}

// Testing of the structs
var j JiraConfig = JiraConfig{
	OnPrem:   false,
	AuthType: "PAT",
}

var rootConfig RootConfig = RootConfig{
	IssueSystem: "jira",
}

func Test() {
	rootConfig.JiraConfigObj.URI = "http://vader.lan:8080"
	rootConfig.JiraConfigObj.SeverityFields.High.Name = "Catastrophic"
	rootConfig.JiraConfigObj.SeverityFields.High.ID = "customfield_100033"
}
