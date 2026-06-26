package services

import (
	"fmt"
	"github-security-scanner/models"
	"os"
	"sort"
	"sync"
	"time"
)
	var LastScan models.ScanResponse
func RunScan(repoURL string) (*models.ScanResponse, error) {
	startTime := time.Now()
	repoPath, err := CloneRepo(repoURL)

	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(repoPath)

	var semgrepVulns []models.Vulnerability
	var gitleaksVuls []models.Vulnerability

	var semgrepErr error
	var gitleaksErr error

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Starting Semgrep")
	fmt.Println("Starting Gitleaks")
	go func() {
		defer wg.Done()
		fmt.Println("semgrep running..")
		semgrepVulns, semgrepErr = RunSemgrep(repoPath)
		fmt.Println("semgrep finished..")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("gitleaks running..")
		gitleaksVuls, gitleaksErr = RunGitleaks(repoPath)
		fmt.Println("gitleaks finished..")
	}()

	wg.Wait()
	if semgrepErr != nil {
		return nil, semgrepErr
	}
	if gitleaksErr != nil {
		return nil, gitleaksErr
	}
	allVulns := append(semgrepVulns, gitleaksVuls...)
	summaryByType := make(map[string]int)
	for _, vuln := range allVulns {
		summaryByType[vuln.Type]++
	}
	errorCount := 0
	warningCount := 0
	infoCount := 0
    
	for _, vuln := range allVulns {
		switch vuln.Severity {
		case "ERROR":
			errorCount++
		case "WARNING":
			warningCount++
		case "INFO":
			infoCount++
		}
	}

	riskScore := (errorCount * 10) + (warningCount * 5) + (infoCount * 1)

	riskLevel := "LOW"
	switch {
	case riskScore > 100:
		riskLevel = "CRITICAL"
	case riskScore > 50:
		riskLevel = "HIGH"
	case riskScore > 20:
		riskLevel = "MEDIUM"
	}
	scanDuration := time.Since(startTime).String()
	severityRank := map[string]int{
		"ERROR":   3,
		"WARNING": 2,
		"INFO":    1,
	}
	sort.Slice(allVulns, func(i, j int) bool {
		return severityRank[allVulns[i].Severity] >
			severityRank[allVulns[j].Severity]
	})

	response := models.ScanResponse{
		RepoURL:              repoURL,
		TotalVulnerabilities: len(allVulns),
		ErrorCount:           errorCount,
		WarningCount:         warningCount,
		InfoCount:            infoCount,
		RiskScore:            riskScore,
		RiskLevel:            riskLevel,
		ScanDuration:         scanDuration,
		SummaryByType:        summaryByType,
		Vulnerabilities:      allVulns,
	}
	LastScan = response
	return &response, nil
}
