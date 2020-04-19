package event

import (
	"encoding/json"

	"app/lib"
	"app/obj"

	log "github.com/sirupsen/logrus"
)

// Register == connect
func Register(c *Client) {
	log.WithField("id", c.ID).Info("User Connected")
}

// Event Catch Message
func Event(c *Client, message []byte) {
	//fmt.Println(string(message))

	var objmap map[string]interface{}

	if err := json.Unmarshal(message, &objmap); err != nil {
		log.WithFields(log.Fields{
			"id":    c.ID,
			"error": err,
		}).Error("JSON Unmarshal Error")
		return
	}

	event, ok := objmap["event"].(string)
	if !ok {
		return
	}

	//fmt.Println(objmap)

	//fmt.Println(event)

	switch event {
	case "init":
		if _, ok := Sockets[c.ID]; ok {
			log.WithField("id", c.ID).Warn("Prevent Login")
			return
		}

		name, ok := objmap["data"].(string)
		if !ok {
			return
		}

		if len(name) > 15 {
			name = ""
		}

		Sockets[c.ID] = c
		obj.Users[c.ID] = obj.NewPlayer(c.ID)
		/*
			var gunList []obj.Gun
			for i := -math.Pi / 2; i <= math.Pi; i += math.Pi / 2 {
				gunList = append(gunList, *obj.NewGun(map[string]interface{}{
					"dir": i,
				}, obj.Users[c.ID].ControlObject, *obj.NewObject(map[string]interface{}{
					"speed":  1,
					"h":      1,
					"damage": 0.65,
					"r":      1,
					"stance": 0.1,
				}, []obj.Gun{*obj.NewGun(map[string]interface{}{
					"delaytime": 3.,
				}, nil, obj.Object{})}, obj.DefaultBulletTick, obj.DefaultCollision, nil, nil)))
			}
			for i := -math.Pi / 4 * 3; i <= math.Pi; i += math.Pi / 2 {
				gunList = append(gunList, *obj.NewGun(map[string]interface{}{
					"dir":      i,
					"waittime": 0.5,
				}, obj.Users[c.ID].ControlObject, *obj.NewObject(map[string]interface{}{
					"speed":  1,
					"h":      1,
					"damage": 0.65,
					"r":      1,
					"stance": 0.1,
				}, []obj.Gun{*obj.NewGun(map[string]interface{}{
					"delaytime": 3.,
				}, nil, obj.Object{})}, obj.DefaultBulletTick, obj.DefaultCollision, nil, nil)))
			}
			obj.Users[c.ID].ControlObject = obj.NewObject(map[string]interface{}{
				"type":     "OctoTank",
				"team":     string(c.ID),
				"name":     name,
				"x":        lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
				"y":        lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
				"h":        50,
				"mh":       50,
				"damage":   20,
				"level":    45,
				"stats":    [8]int{0, 0, 0, 7, 7, 7, 7, 5},
				"maxStats": [8]int{7, 7, 7, 7, 7, 7, 7, 7},
			}, gunList, obj.TankTick, obj.DefaultCollision, nil, nil)
		*/
		obj.Users[c.ID].ControlObject = obj.NewObject(map[string]interface{}{
			"type":     "Necromanser",
			"team":     string(c.ID),
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
		}, nil, obj.Object{})}, obj.TankTick, obj.DefaultCollision, obj.NecroKillEvent, nil)
		obj.Users[c.ID].ControlObject.SetController(obj.Users[c.ID])
		obj.Objects = append(obj.Objects, obj.Users[c.ID].ControlObject)

	case "input":
		if m, ok := objmap["data"].(map[string]interface{}); ok {
			if t, ok := m["type"].(string); ok {
				if u, ok := obj.Users[c.ID]; ok {
					switch t {
					case "moveVector":
						if value, ok := m["value"].(float64); ok {
							u.IsMove = value != 9
							u.MoveDir = value
						}
					case "mouseMove":
						pos, ok := m["value"].(map[string]interface{})
						if !ok {
							return
						}
						if obj := u.ControlObject; obj != nil {
							u.SetMousePoint(pos["x"].(float64)-obj.X, pos["y"].(float64)-obj.Y)
						} else {
							u.SetMousePoint(pos["x"].(float64), pos["y"].(float64))
						}
					case "mouseLeft":
						if value, ok := m["value"].(bool); ok {
							u.Ml = value
						} else {
							u.Ml = false
						}
					case "mouseRight":
						if value, ok := m["value"].(bool); ok {
							u.Mr = value
						} else {
							u.Mr = false
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
}

// UnRegister == close
func UnRegister(c *Client) {
	if _, ok := Sockets[c.ID]; !ok {
		log.WithField("id", c.ID).Warn("Prevent Disconnect")
		return
	}

	obj.Users[c.ID].ControlObject.Controller = nil

	delete(Sockets, c.ID)
	delete(obj.Users, c.ID)

	log.WithField("id", c.ID).Info("User Disconnect")
}
