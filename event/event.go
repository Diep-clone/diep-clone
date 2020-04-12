package event

import "fmt"

// Event Catch Message
func Event(c *Client, message []byte) {
	fmt.Println(string(message))
}

// Close event
func Close() {

}