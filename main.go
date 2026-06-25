package main

import (
	"github-security-scanner/services"

	"github.com/gin-gonic/gin"
)

type ScanReq struct {
	RepoURL string `json:"repo_url"`
}

func main() {

	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
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
	r.Run(":8080")
}
