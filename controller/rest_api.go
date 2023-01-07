package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SpawnServer() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/prompts", getAllPrompts)
	router.POST("/prompts", createPrompt)
	router.PATCH("/prompts/toggle", togglePrompt)
	router.GET("/stats/unsafe/:query", getUnsafeStats)
	router.GET("/stats/reports/:query", readReports)

	router.POST("/hadoop/run/:query", runHadoopJob)
	router.GET("/hadoop/finished/:query/:jobName", notifyFinishedJob)

	router.Run("0.0.0.0:5321")
}

func getAllPrompts(c *gin.Context) {
	prompts := mongoCtx.GetAllPrompts()
	c.IndentedJSON(http.StatusOK, prompts)
}

func getUnsafeStats(c *gin.Context) {
	query := c.Param("query")
	stats := mongoCtx.ReadAllUnsafeStats(query)
	c.IndentedJSON(http.StatusOK, stats)
}

func togglePrompt(c *gin.Context) {
	var prompt Prompt
	err := c.BindJSON(&prompt)
	if err != nil {
		log.Println(err)
	}
	mongoCtx.TogglePromptStatus(prompt.Query)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Update successful"})
}

func createPrompt(c *gin.Context) {
	var prompt Prompt
	err := c.BindJSON(&prompt)
	if err != nil {
		panic(err)
	}
	mongoCtx.InsertPrompt(prompt)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Create successful"})
}

func runHadoopJob(c *gin.Context) {
	query := c.Param("query")
	uuid := strings.Replace(uuid.New().String(), "-", "", -1)
	ExecHadoopJob(query, uuid)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Schedule successful"})
}

func notifyFinishedJob(c *gin.Context) {
	query := c.Param("query")
	jobName := c.Param("jobName")
	content := hdfsConn.ReadFromFile(query, jobName)
	report := Report{Query: query, Data: content, Id: jobName, Timestamp: time.Now().UnixNano()}

	mongoCtx.InsertReport(report)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successful notification"})
}

func readReports(c *gin.Context) {
	query := c.Param("query")
	reports := mongoCtx.ReadAllReports(query)

	c.IndentedJSON(http.StatusOK, reports)
}
