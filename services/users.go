package users

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var userList []User
var lastUserID int

func CreateUser(c *gin.Context) {
	var newUserInput User

	// Bind the JSON data from the request body to newUserInput
	if err := c.BindJSON(&newUserInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	// Validate user input , the username cannot be empty
	if newUserInput.Name == "" || strings.TrimSpace(newUserInput.Name) == "" || newUserInput.Age <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and age are required and must be valid"})
		return
	}

	// Generate a unique user ID for the user created above
	newUserID := lastUserID + 1
	lastUserID = newUserID

	// Create a new user with a unique ID generated above
	newUser := User{
		ID:   newUserID,
		Name: newUserInput.Name,
		Age:  newUserInput.Age,
	}

	// Append the new user to the userList , this contains all the list of users , will help in retrieving all the users
	userList = append(userList, newUser)

	// Respond with the newly created user's details after adding the newly created user
	c.JSON(http.StatusCreated, newUser)
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userList)
}
