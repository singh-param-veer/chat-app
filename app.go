package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Age  float64 `json:"age"`
}

var users = []user{
	{ID: "1", Name: "Shashank Rai", Age: 27},
	{ID: "2", Name: "Saurabh Kumar", Age: 28},
	{ID: "3", Name: "Prashant Sharma", Age: 27},
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)

	router.Run("localhost:8080")
}
