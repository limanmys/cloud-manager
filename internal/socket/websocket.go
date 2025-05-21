package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var mut map[string]*sync.Mutex = make(map[string]*sync.Mutex)

type miniClient map[string]*websocket.Conn // Modified type

type ClientObject struct {
	MachineId string
	conn      *websocket.Conn
}

type ResponseObject struct {
	MSG  string
	FROM ClientObject
}

var clients = make(miniClient)
var register = make(chan ClientObject)
var response = make(chan ResponseObject)
var unregister = make(chan ClientObject)

var responses = make(map[string]string)

func removeClient(machine_id string) {
	if conn, ok := clients[machine_id]; ok { // Check if client exists
		delete(clients, machine_id)
		conn.Close() // Close the connection before potentially removing the organization map

	}
}

func WebSocketHandler() {
	for {
		select {
		case client := <-register:
			// Pre-initialize organization map if it doesn't exist
			if clients[client.MachineId] == nil {
				clients[client.MachineId] = &websocket.Conn{}
			}

			if mut[client.MachineId] == nil {
				mut[client.MachineId] = &sync.Mutex{}
			}

			clients[client.MachineId] = client.conn
			log.Println("client registered:", client.MachineId)

		case message := <-response:
			mut[message.FROM.MachineId].Lock()
			responses[message.FROM.MachineId] = message.MSG
			mut[message.FROM.MachineId].Unlock()

		case client := <-unregister:
			removeClient(client.MachineId) // Update client removal
			log.Println("client unregistered:", client.MachineId)
		}
	}
}

func Init(app *fiber.App) {
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		clientObj := ClientObject{
			MachineId: c.Locals("machine_id").(string),
			conn:      c,
		}
		defer func() {
			unregister <- clientObj
			c.Close()
		}()
		// Register the client
		register <- clientObj

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				response <- ResponseObject{
					MSG:  string(message),
					FROM: clientObj,
				}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))
}

func Send(machine_id string, command string, data string) (string, error) {

	if clients[machine_id] == nil {
		return "", fmt.Errorf("device is not connected")
	}

	responses[machine_id] = ""
	payload := make(map[string]string)
	payload["command"] = command
	payload["data"] = data

	bt, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	fmt.Println("sending:", string(bt))

	err = clients[machine_id].WriteMessage(websocket.TextMessage, bt)
	if err != nil {
		removeClient(machine_id) // Update client removal
		clients[machine_id].WriteMessage(websocket.CloseMessage, []byte{})
		clients[machine_id].Close()
		return "", err
	}

	timeout := 0
	for {
		if timeout > 10 {
			return "", fmt.Errorf("socket timeout")
		}

		mut[machine_id].Lock()

		if responses[machine_id] == "" {
			time.Sleep(500 * time.Millisecond)
		} else {
			mut[machine_id].Unlock()
			break
		}
		mut[machine_id].Unlock()
		timeout++

	}

	return responses[machine_id], nil
}
