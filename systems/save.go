package systems

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type SaveData struct {
	Coins     int            `json:"coins"`
	BestScore int            `json:"best_score"`
	Upgrades  map[string]int `json:"upgrades"`
}

func LoadSave() SaveData {
	data := SaveData{
		Upgrades: make(map[string]int),
	}

	file, err := ioutil.ReadFile("save.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Return default save
			return data
		}
		return data
	}

	json.Unmarshal(file, &data)
	if data.Upgrades == nil {
		data.Upgrades = make(map[string]int)
	}
	return data
}

func SaveGame(data SaveData) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("save.json", file, 0644)
}
