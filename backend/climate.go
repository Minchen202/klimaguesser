package main

import (
	"encoding/json"
	"os"
)

func loadClimateData(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &climateData)
}

func closestClimate(lat, lng float64) Climate {
	var closest Climate
	closestDistance := -1.0

	for _, c := range climateData {
		d := calculateDistance(lat, lng, c.lat(), c.lng())
		if closestDistance < 0 || d < closestDistance {
			closestDistance = d
			closest = c
		}
	}
	return closest
}
