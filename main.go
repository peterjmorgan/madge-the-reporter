package main

import "github.com/peterjmorgan/madge-the-reporter/cmd"

///*
//Reqes:
//- JIRA project
//- Phylum project ID (can read from .phylum_project file)
//*/
//
//type PhylumProjectFile struct {
//	ID        string `yaml:"id"`
//	Name      string `yaml:"name"`
//	CreatedAt string `yaml:"created_at"`
//}
//
//// Read Phylum Project ID from .phylum_project file
//// This might not be necessary
//func readPhylumProjectFile() (string, error) {
//	_, err := os.Stat(".phylum_project")
//	if err != nil {
//		log.Errorf("failed to stat .phylum_project: %v\n", err)
//		return "", err
//	}
//
//	phylumFileData, err := os.ReadFile(".phylum_project")
//	if err != nil {
//		log.Errorf("failed to read .phylum_project: %v\n", err)
//		return "", err
//	}
//
//	phylumConfig := new(structs.PhylumProject)
//	err = yaml.Unmarshal(phylumFileData, phylumConfig)
//	if err != nil {
//		log.Errorf("failed to unmarshall config data: %v\n", err)
//		return "", nil
//	}
//
//	return phylumConfig.ID, nil
//}

func main() {
	cmd.Execute()
}
