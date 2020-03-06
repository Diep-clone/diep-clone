package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Setting is
type Setting struct {
	GameMode  string `json:"gameMode"`
	MaxPlayer int    `json:"maxPlayer"`
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
