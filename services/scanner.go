package services

import (
	"fmt"
	"github-security-scanner/models"
	"os"
	"path/filepath"
	"strings"
)

func ScanFiles() ([]models.Vulnerability, error) {
	vulnerabilities := []models.Vulnerability{}

	err := filepath.Walk("./repos/project",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println("Scanning :", path)
			data, err := os.ReadFile(path)
			if err == nil {
				content := strings.ToLower(string(data))
				if strings.Contains(content, "password") ||
					strings.Contains(content, "secret") ||
					strings.Contains(content, "api-keys") {
					vulnerabilities = append(vulnerabilities,
						models.Vulnerability{
							Type:        "Hardcoded Secret",
							File:        path,
							Description: "Possible Secret Found",
						},
					)
					fmt.Println("Potential secret found:", path)
				}
			}
			return nil
		},
	)
	return vulnerabilities, err
}
