package event

import (
	"encoding/json"
	"fmt"
	"log"
)

// Register == connect
func Register(c *Client) {
	fmt.Println("hello")
}

// Event Catch Message
func Event(c *Client, message []byte) {
	fmt.Println(string(message))

	var objmap map[string]interface{}

	err := json.Unmarshal(message, &objmap)
	if err != nil {
		log.Println(err)
	}

	event, ok := objmap["event"].(string)
	if !ok {
		return
	}

	switch event {
	case "init":
	case "input":

	}
}

// UnRegister == close
func UnRegister(c *Client) {
	fmt.Println("bye")
}
