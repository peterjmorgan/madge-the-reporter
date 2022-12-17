package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	phylum "github.com/peterjmorgan/go-phylum"
)

/*
Reqes:
- JIRA project
- Phylum project ID (can read from .phylum_project file)
 */

type PhylumProjectFile struct {
	ID string `yaml:"id"`
	Name string `yaml:"name"`
	CreatedAt string `yaml:"created_at"`
}


func readPhylumProjectFile() (string, error) {
	_, err := os.Stat(".phylum_project")
	if err != nil {
		log.Errorf("failed to stat .phylum_project: %v\n", err)
		return "", err
	}
	phylumFileData, err := os.ReadFile(".phylum_project")
	if err != nil {
		log.Errorf("failed to read .phylum_project: %v\n", err)
		return "", err
	}
	yaml.


}

func main() {
	projectId := os.Args[1]
	if len(projectId) < 32 {
		fmt.Printf("failed to read argument. Please provide a J")
	}
}