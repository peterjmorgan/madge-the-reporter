package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/peterjmorgan/madge-the-reporter/structs"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ReadConfigFile(testConfigData *structs.ConfigFile) (*structs.RootConfig, error) {
	var filename string = "madge_config.yaml"

	// For testing
	v := reflect.ValueOf(testConfigData)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		if testConfigData.Filename != "" {
			filename = testConfigData.Filename
		}
	}

	if _, err := os.Stat(filename); err == nil {
		// exists
		fileData, err1 := os.ReadFile(filename)
		if err1 != nil {
			fmt.Printf("Failed to read madge_config.yaml: %v\n", err1)
			return nil, fmt.Errorf("failed to read file")
		}

		configData := new(structs.RootConfig)
		err2 := yaml.Unmarshal(fileData, configData)
		if err2 != nil {
			fmt.Printf("Failed to unmarshall config data: %v\n", err2)
			return nil, fmt.Errorf("failed to unmarshall config data")
		}

		return configData, nil
	}
	return nil, fmt.Errorf("config file not found")
}

func PromptForString(message string, lenRequirement int) (string, error) {
	prompt := promptui.Prompt{
		Label: message,
		Validate: func(input string) error {
			strLen := len(input)
			if strLen != lenRequirement && lenRequirement != -1 {
				return errors.New("invalid length")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("PromptForString: %v failed:%v\n", message, err)
		return "", err
	}
	return result, nil
}

func PhylumGetAuthToken() (string, error) {
	var retStr string
	var stdErrBytes bytes.Buffer

	var authTokenArgs = []string{"auth", "token"}
	authTokenCmd := exec.Command("phylum", authTokenArgs...)
	authTokenCmd.Stderr = &stdErrBytes

	retBytes, err := authTokenCmd.Output()
	if err != nil {
		log.Errorf("Failed to exec 'phylum auth token': %v\n", err)
		log.Errorf(stdErrBytes.String())
		return "", err
	}
	stdErrString := stdErrBytes.String()
	_ = stdErrString // prob will need this later

	retStr = string(retBytes)
	retStr = strings.Trim(retStr, "\n")

	return retStr, nil
}
