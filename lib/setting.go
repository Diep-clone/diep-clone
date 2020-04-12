package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var GameSetting Setting = ReadSetting()

// Setting is
type Setting struct {
	GameMode  string  `json:"gameMode"`
	GameSpeed float64 `json:"gameSpeed"`
	MaxPlayer int     `json:"maxPlayer"`
	MaxShape  int     `json:"maxShape"`
	MapSize   struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"mapSize"`
}

// ReadSetting returns Setting
func ReadSetting() Setting {
	settingFile, osErr := os.Open("./config/setting.json")
	if osErr != nil {
		log.Fatal(osErr)
	}

	defer settingFile.Close()

	byteVal, _ := ioutil.ReadAll(settingFile)

	var setting Setting
	if err := json.Unmarshal(byteVal, &setting); err != nil {
		log.Fatal(err)
	}

	return setting
}
