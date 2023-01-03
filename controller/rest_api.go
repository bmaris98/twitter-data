package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SpawnServer() {
	router := gin.Default()
	router.GET("/prompts", getAllPrompts)
	router.POST("/prompts", createPrompt)
	router.PATCH("/prompts/toggle", togglePrompt)
	router.GET("/stats/unsafe/:query", getUnsafeStats)

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
