package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	model "web-service/model"
	users "web-service/services"

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

var activeSocketUsers map[int]*websocket.Conn = make(map[int]*websocket.Conn)

func startSocketMessaging(conn *websocket.Conn) {
	// Continuosly read and write message
	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}
		fmt.Println(string(rawMessage))
		body := model.ParseToModel(rawMessage)
		fmt.Printf("body: %v\n", body)
		fmt.Printf("body: %v\n", body.GetCmd())
		if body.GetCmd() == "add" {
			todoList = append(todoList, body.GetTodo())
		} else if body.GetCmd() == "done" {
			updateTodoList(body.GetTodo())
		}

		messageRedirecting(body.Receiver, body)

		output := "Ongoing Tasks: \n"
		if len(todoList) == 0 {
			output += "Empty"
		} else {
			for _, todo := range todoList {
				output += "\n - " + todo + "\n"
			}
		}
		sendMessage(conn, output)
	}
}

func sendMessage(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write failed:", err)
		// conn.Close()
	}
}

func messageRedirecting(receiver int, message model.Message) {
	if activeSocketUsers[receiver] != nil {
		sendMessage(activeSocketUsers[receiver], string(model.TransformToJson(message)))
	}
}

func main() {
	router := gin.Default()
	//POST endpoint to create new users
	router.POST("/users", users.CreateUser)
	// GET endpoint to retrieve all users
	router.GET("/users", users.GetUsers)
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		query_params := c.Request.URL.Query()
		user_id, err := strconv.Atoi(query_params.Get("user_id"))

		if err != nil {
			panic("Please set user id as query params in the websocket request")
		}

		if activeSocketUsers[user_id] != nil {
			return
		}

		activeSocketUsers[user_id] = conn
		conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!. You are now connected"))
		startSocketMessaging(conn)
	})

	router.Run("localhost:8080")
}
