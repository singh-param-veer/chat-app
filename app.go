package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var todoList []string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getActionAndTask(input string) (string, string) {
	inputArr := strings.Split(input, " ")
	return inputArr[0], inputArr[1]
}

func getMessage(input string) string {
	inputArr := strings.Split(input, " ")
	var result string
	for i := 1; i < len(inputArr); i++ {
		result += inputArr[i]
	}
	return result
}

func updateTodoList(input string) {
	tmpList := todoList
	todoList = []string{}
	for _, val := range tmpList {
		if val == input {
			continue
		}
		todoList = append(todoList, val)
	}
}

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

func createWebsocket(c *gin.Context) {
	c.IndentedJSON(http.StatusCreated, "Created Websocket")
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!. You are now connected"))

		// Continuosly read and write message
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}
			input := string(message)
			cmd, msg := getActionAndTask(input)
			if cmd == "add" {
				todoList = append(todoList, msg)
			} else if cmd == "done" {
				updateTodoList(msg)
			}
			output := "Ongoing Tasks: \n"
			if len(todoList) == 0 {
				output += "Empty"
			} else {
				for _, todo := range todoList {
					output += "\n - " + todo + "\n"
				}
			}
			message = []byte(output)
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write failed:", err)
				break
			}
		}
	})

	router.Run("localhost:8080")
}
