package event

import (
	"encoding/json"
	"fmt"
)

// Event Catch Message
func Event(c *Client, message []byte) {
	fmt.Println(string(message))

	var objmap map[string]interface{}
	_ = json.Unmarshal(message, &objmap)
	event := objmap["event"].(string)

	switch event {
	case "init":
	case "input":

	}
}

// Close event
func Close() {

}
