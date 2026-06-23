package main

import (
	"github-security-scanner/models"
	"github-security-scanner/services"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type ScanReq struct {
	RepoURL string `json:"repo_url"`
}

func main() {
	r := gin.Default()

	r.POST("/scan", func(c *gin.Context) {
		startTime := time.Now()
		var req ScanReq

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid Request",
			})
			return
		}
		err := services.CloneRepo(req.RepoURL)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		// vulnerabilities, err := services.ScanFiles()

		semgrepVulns, err := services.RunSemgrep()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		// vulnerabilities = append(vulnerabilities, semgrepVulns...)
		// println(len(semgrepVulns))

		// println(string(output))

		gitleaksVuls, err := services.RunGitleaks()
		if err != nil {
			c.JSON(500, gin.H{
				"err": err.Error(),
			})
			return
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
			RepoURL:              req.RepoURL,
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
		c.JSON(200, response)
	})
	r.Run(":8080")
}
