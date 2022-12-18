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

	router.Run("localhost:8080")
}

func getAllPrompts(c *gin.Context) {
	prompts := mongoCtx.GetAllPrompts()
	c.IndentedJSON(http.StatusOK, prompts)
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
