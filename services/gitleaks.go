package services

import (
	"encoding/json"
	"github-security-scanner/models"
	"os"
	"os/exec"
)

func RunGitleaks() ([]models.Vulnerability, error) {
	cmd := exec.Command(
		"gitleaks",
		"detect",
		"--source=./repos/project",
		"--report-format=json",
		"--report-path=gitleaks.json",
	)
	err := cmd.Run()
	_ = err
	data, err := os.ReadFile("gitleaks.json")
	if err != nil {
		return nil, err
	}
	var findings []models.GitleaksFinding
	json.Unmarshal(data, &findings)
	if err != nil {
		return nil, err
	}
	var vulnerabilities []models.Vulnerability
	for _, finding := range findings {
		vulnerabilities = append(vulnerabilities,
			models.Vulnerability{
				Type:        finding.RuleID,
				File:        finding.File,
				Line:        finding.StartLine,
				Description: finding.Description,
				Severity:    "ERROR",
			},
		)
	}

	return vulnerabilities, nil
}
