package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
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
		log.WithError(osErr).Error("Read Setting.json Error")
	}

	defer settingFile.Close()

	byteVal, _ := ioutil.ReadAll(settingFile)

	var setting Setting
	if err := json.Unmarshal(byteVal, &setting); err != nil {
		log.WithError(err).Error("Setting.json Unmarshal Error")
	}

	return setting
}
