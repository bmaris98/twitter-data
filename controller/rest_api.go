package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SpawnServer() {
	router := gin.Default()

	// config := cors.Config{
	// 	AllowAllOrigins: true,
	// 	AllowMethods:    []string{"GET", "PATCH", "POST", "OPTIONS"},
	// }
	// config.AddAllowHeaders("*")

	//router.Use(cors.New(config))

	router.Use(cors.Default())

	router.GET("/prompts", getAllPrompts)
	router.POST("/prompts", createPrompt)
	router.PATCH("/prompts/toggle", togglePrompt)
	router.GET("/stats/unsafe/:query", getUnsafeStats)

	router.POST("/hadoop/run/:query", runHadoopJob)

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
