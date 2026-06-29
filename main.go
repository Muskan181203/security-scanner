package main

import (
	"github-security-scanner/services"

	"github.com/gin-gonic/gin"
)

type ScanReq struct {
	RepoURL string `json:"repo_url"`
}
type EmailReq struct {
	Email string `json:"email"`
}

func main() {

	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.GET("/report", func(c *gin.Context) {

		filePath, err := services.GenerateReport(
			*services.LastScan,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.FileAttachment(
			filePath,
			"security-report.html",
		)
	})
	r.GET("/report/pdf", func(c *gin.Context) {

		file, err := services.GeneratePDFReport(*services.LastScan)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.FileAttachment(file, "security-report.pdf")
	})
	r.POST("/scan", func(c *gin.Context) {

		var req ScanReq

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid Request",
			})
			return
		}

		response, err := services.RunScan(req.RepoURL)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, response)

	})

	r.POST("/send-email", func(c *gin.Context) {

		var req EmailReq

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid Request",
			})
			return
		}

		// Make sure a scan has already been performed
		if services.LastScan == nil {
			c.JSON(400, gin.H{
				"error": "Please scan a repository first.",
			})
			return
		}

		// Generate the latest PDF
		pdfFile, err := services.GeneratePDFReport(*services.LastScan)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Send email with PDF attachment
		err = services.SendEmail(req.Email, pdfFile)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Email sent successfully!",
		})
	})
	r.Run(":8080")
}
