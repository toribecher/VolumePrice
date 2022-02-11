package vwap

import (
	"VolumePrice/helper"
	"fmt"
	"strconv"
)

const maxLimit = 200

func GetVWap(input chan helper.Match) {
	var pairInfo helper.PairInfo
	pairInfo.Matches = make(map[string][]helper.Match)
	pairInfo.TotalShares = make(map[string]float64)
	pairInfo.TotalSpent = make(map[string]float64)
	pairInfo.VolumeWeightedAveragePrice = make(map[string]float64)
	for match := range input {
		if match.Size == "" {
			continue
		}
		var size float64
		var price float64
		var err error
		size, err = strconv.ParseFloat(match.Size, 64)
		price, err = strconv.ParseFloat(match.Price, 64)
		if err != nil {
			fmt.Printf("string conversion error %s", err)
			return
		}
		if len(pairInfo.Matches[match.ProductId]) > maxLimit {
			pairInfo.Matches[match.ProductId] = append(pairInfo.Matches[match.ProductId][:0], pairInfo.Matches[match.ProductId][1:]...)
		}
		totalSpent := size * price
		volumeWeightedAverage := calculateVolumeWeightedAveragePrice(pairInfo.TotalSpent[match.ProductId]+totalSpent, pairInfo.TotalShares[match.ProductId]+size)
		pairInfo.VolumeWeightedAveragePrice[match.ProductId] = volumeWeightedAverage
		pairInfo.TotalSpent[match.ProductId] = totalSpent + size
		pairInfo.Matches[match.ProductId] = append(pairInfo.Matches[match.ProductId], match)
		fmt.Println(pairInfo.VolumeWeightedAveragePrice)
	}
}

func calculateVolumeWeightedAveragePrice(totalSpent, totalSharesBought float64) float64 {
	return totalSpent / totalSharesBought
}
