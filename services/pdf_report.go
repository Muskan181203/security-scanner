package services

import (
	"fmt"
	"os"

	"github-security-scanner/models"

	"github.com/phpdave11/gofpdf"
)

func GeneratePDFReport(scan models.ScanResponse) (string, error) {

	fileName := "security-report.pdf"

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetTitle("GitHub Security Report", false)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(190, 10, "GitHub Security Report")
	pdf.Ln(15)

	// Repository
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Repository")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(190, 7, scan.RepoURL, "", "", false)
	pdf.Ln(3)

	// Risk Summary
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Risk Summary")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(60, 8, fmt.Sprintf("Risk Score: %d", scan.RiskScore))
	pdf.Ln(7)

	pdf.Cell(60, 8, fmt.Sprintf("Risk Level: %s", scan.RiskLevel))
	pdf.Ln(7)

	pdf.Cell(60, 8, fmt.Sprintf("Total Vulnerabilities: %d", scan.TotalVulnerabilities))
	pdf.Ln(7)

	pdf.Cell(60, 8, fmt.Sprintf("Errors: %d", scan.ErrorCount))
	pdf.Ln(7)

	pdf.Cell(60, 8, fmt.Sprintf("Warnings: %d", scan.WarningCount))
	pdf.Ln(7)

	pdf.Cell(60, 8, fmt.Sprintf("Info: %d", scan.InfoCount))
	pdf.Ln(7)

	pdf.Cell(60, 8, fmt.Sprintf("Scan Duration: %s", scan.ScanDuration))
	pdf.Ln(12)

	// Summary by Type
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Summary By Type")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)

	for vulnType, count := range scan.SummaryByType {
		pdf.Cell(190, 7, fmt.Sprintf("%s : %d", vulnType, count))
		pdf.Ln(6)
	}

	pdf.Ln(5)

	// Vulnerabilities
	// -----------------------------
	// Vulnerabilities
	// -----------------------------

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Vulnerabilities")
	pdf.Ln(12)

	for i, v := range scan.Vulnerabilities {

		// Draw a separator line between vulnerabilities
		if i != 0 {
			pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
			pdf.Ln(4)
		}

		// Vulnerability title
		pdf.SetFont("Arial", "B", 13)
		pdf.Cell(190, 8, fmt.Sprintf("Vulnerability #%d", i+1))
		pdf.Ln(8)

		// Severity
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(30, 7, "Severity :")

		switch v.Severity {
		case "ERROR":
			pdf.SetTextColor(220, 53, 69) // Red
		case "WARNING":
			pdf.SetTextColor(253, 126, 20) // Orange
		case "INFO":
			pdf.SetTextColor(13, 110, 253) // Blue
		default:
			pdf.SetTextColor(0, 0, 0)
		}

		pdf.SetFont("Arial", "", 11)
		pdf.Cell(100, 7, v.Severity)
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(7)

		// Type
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(30, 7, "Type :")

		pdf.SetFont("Arial", "", 11)
		pdf.MultiCell(160, 7, v.Type, "", "", false)

		// File
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(30, 7, "File :")

		pdf.SetFont("Arial", "", 11)
		pdf.MultiCell(160, 7, v.File, "", "", false)

		// Line
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(30, 7, "Line :")

		pdf.SetFont("Arial", "", 11)
		pdf.Cell(30, 7, fmt.Sprintf("%d", v.Line))
		pdf.Ln(8)

		// Description
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(30, 7, "Description :")
		pdf.Ln(7)

		pdf.SetFont("Arial", "", 11)
		pdf.MultiCell(190, 7, v.Description, "", "", false)

		pdf.Ln(8)

		// New page if needed
		if pdf.GetY() > 260 {
			pdf.AddPage()
		}
	}

	// Table Header
	pdf.SetFont("Arial", "B", 10)

	pdf.CellFormat(25, 8, "Severity", "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, 8, "Type", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, "File", "1", 0, "C", false, 0, "")
	pdf.CellFormat(15, 8, "Line", "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, 8, "Description", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)

	for _, v := range scan.Vulnerabilities {

		// Change severity color
		switch v.Severity {
		case "ERROR":
			pdf.SetTextColor(220, 53, 69)
		case "WARNING":
			pdf.SetTextColor(253, 126, 20)
		case "INFO":
			pdf.SetTextColor(13, 110, 253)
		default:
			pdf.SetTextColor(0, 0, 0)
		}

		pdf.CellFormat(25, 8, v.Severity, "1", 0, "", false, 0, "")

		pdf.SetTextColor(0, 0, 0)

		pdf.CellFormat(45, 8, v.Type, "1", 0, "", false, 0, "")
		pdf.CellFormat(60, 8, v.File, "1", 0, "", false, 0, "")
		pdf.CellFormat(15, 8, fmt.Sprintf("%d", v.Line), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, 8, v.Description, "1", 1, "", false, 0, "")
	}

	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	_, err = os.Stat(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
