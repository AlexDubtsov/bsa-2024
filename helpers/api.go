package helpers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetJson(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func GetRate(url, pair string) float64 {
	var fetchedStruct RateApiAllRecords
	var currencyRateStr string

	// Open external API
	GetJson(url, &fetchedStruct)

	// Currency pair rate choice
	for _, word := range fetchedStruct.Data {
		if word.Symbol == pair {
			currencyRateStr = word.Value
		}
	}
	// Convert string to float64
	resultRate, _ := strconv.ParseFloat(currencyRateStr, 64)

	return resultRate
}
