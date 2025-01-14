package util

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/open-policy-agent/opa/loader"
)

type Rule struct {
	PublicId string
	Path     string
}

func RetrieveRules(paths []string) ([]Rule, error) {
	fileNames, err := findRegoFiles(paths)
	if err != nil {
		return []Rule{}, err
	}

	rules := []Rule{}
	for _, fileName := range fileNames {
		publicId, err := getPublicIdFromFile(fileName)
		if err != nil {
			return []Rule{}, err
		}
		if publicId != "" {
			rules = append(rules, Rule{
				PublicId: publicId,
				Path:     fileName,
			})
		}
	}
	return rules, nil
}

func findRegoFiles(paths []string) ([]string, error) {
	fileNames := []string{}

	result, err := loader.AllRegos(paths)
	if err != nil {
		return fileNames, err
	}

	for _, module := range result.Modules {
		fileName := filepath.Clean(module.Name)
		if !strings.Contains(fileName, "_test.rego") {
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames, nil
}

func getPublicIdFromFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	publicId := extractPublicIdFromRego(string(data))
	return publicId, nil
}

func extractPublicIdFromRego(rego string) string {
	re := regexp.MustCompile("\"publicId\"\\s*:\\s*\"(.*?)\"")
	match := re.FindStringSubmatch(rego)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}
