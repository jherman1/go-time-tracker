package tracker

import (
	"encoding/json"
	"os"
)

const dataFile = "data/tracker_data.json"

func LoadTracker() (*Tracker, error) {
	f, err := os.Open(dataFile)
	if err != nil {
		return &Tracker{}, nil //return empty if the file doesn't exist
	}
	defer f.Close()

	var t Tracker
	err = json.NewDecoder(f).Decode(&t)
	return &t, err
}

func SaveTracker(t *Tracker) error {
	os.MkdirAll("data", os.ModePerm) // <-- Ensures the data directory exists

	f, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(t)
}

func DeleteTrackingFile() error {
	return os.Remove(dataFile)
}
