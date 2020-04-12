package event

import (
	"encoding/json"
	"fmt"
	"log"

	"app/lib"
	"app/obj"
)

// Register == connect
func Register(c *Client) {
	log.Println("INFO > " + string(c.ID) + " Connected")
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
		if _, ok := Sockets[c.ID]; ok {
			log.Println("INFO > Prevent " + string(c.ID) + " Login")
			return
		}
		name, ok := objmap["name"].(string)
		if !ok {
			return
		}

		if len(name) > 15 {
			name = ""
		}

		Sockets[c.ID] = c
		obj.Users[c.ID] = obj.NewPlayer(c.ID)
		obj.Users[c.ID].ControlObject = obj.NewObject(map[string]interface{}{
			"type":     "Necromanser",
			"team":     c.ID,
			"name":     name,
			"x":        lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
			"y":        lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
			"h":        50,
			"mh":       50,
			"damage":   20,
			"level":    45,
			"stats":    [8]int{0, 0, 0, 7, 7, 7, 7, 5},
			"maxStats": [8]int{7, 7, 7, 7, 7, 7, 7, 7},
			"sight":    1.11,
		}, []obj.Gun{*obj.NewGun(map[string]interface{}{
			"owner": nil,
			"limit": 0,
		})}, obj.TankTick, obj.DefaultCollision, obj.NecroKillEvent, nil)
		obj.Users[c.ID].ControlObject.SetController(obj.Users[c.ID])
		obj.Objects = append(obj.Objects, obj.Users[c.ID].ControlObject)

	case "input":
		if t, ok := objmap["type"].(string); ok {
			if u, ok := obj.Users[c.ID]; ok {
				switch t {
				case "moveVector":
					if value, ok := objmap["value"].(float64); ok {
						u.IsMove = value != 9
						u.MoveDir = value
					}
				case "mouseMove":
					x, ok := objmap["x"].(float64)
					if !ok {
						return
					}
					y, ok := objmap["y"].(float64)
					if !ok {
						return
					}
					if obj := u.ControlObject; obj != nil {
						u.SetMousePoint(x-obj.X, y-obj.Y)
					} else {
						u.SetMousePoint(x, y)
					}
				case "mouseLeft":
					if value, ok := objmap["value"].(bool); ok {
						u.Ml = value
					}
				case "mouseRight":
					if value, ok := objmap["value"].(bool); ok {
						u.Mr = value
					}
				case "suicide":
					if obj := u.ControlObject; obj != nil {
						obj.H = 0
					}
				}
			}
		}
	}
}

// UnRegister == close
func UnRegister(c *Client) {
	if _, ok := Sockets[c.ID]; !ok {
		log.Println("INFO > Prevent " + string(c.ID) + " Disconnect")
		return
	}

	obj.Users[c.ID].ControlObject.Controller = nil

	delete(Sockets, c.ID)
	delete(obj.Users, c.ID)

	log.Println("INFO > " + string(c.ID) + " disconnected")
}
