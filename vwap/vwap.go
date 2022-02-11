package vwap

import (
	"VolumePrice/helper"
	"fmt"
	"log"
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
		printedPairInfo := doCalculations(match, pairInfo)
		log.Println(printedPairInfo.VolumeWeightedAveragePrice)
	}
}

func doCalculations(match helper.Match, pairInfo helper.PairInfo) helper.PairInfo {
	if match.Size == "" {
		return helper.PairInfo{}
	}
	var size float64
	var price float64
	var err error
	size, err = strconv.ParseFloat(match.Size, 64)
	price, err = strconv.ParseFloat(match.Price, 64)
	if err != nil {
		fmt.Printf("string conversion error %s", err)
		return helper.PairInfo{}
	}
	removePairCheck(pairInfo, match.ProductId, maxLimit, match)

	totalSpent := size * price
	volumeWeightedAverage := calculateVolumeWeightedAveragePrice(pairInfo.TotalSpent[match.ProductId]+totalSpent, pairInfo.TotalShares[match.ProductId]+size)
	pairInfo.VolumeWeightedAveragePrice[match.ProductId] = volumeWeightedAverage
	pairInfo.TotalSpent[match.ProductId] = pairInfo.TotalSpent[match.ProductId] + totalSpent
	pairInfo.TotalShares[match.ProductId] = pairInfo.TotalShares[match.ProductId] + size
	pairInfo.Matches[match.ProductId] = append(pairInfo.Matches[match.ProductId], match)
	return pairInfo
}

func removePairCheck(pairInfo helper.PairInfo, productId string, maxLimit int, match helper.Match) {
	if len(pairInfo.Matches[match.ProductId]) > maxLimit {
		var size float64
		var price float64
		var err error
		size, err = strconv.ParseFloat(match.Size, 64)
		price, err = strconv.ParseFloat(match.Price, 64)
		if err != nil {
			fmt.Println("parsing error")
			return
		}
		totalSpent := price * size
		pairInfo.TotalSpent[productId] = pairInfo.TotalSpent[productId] - totalSpent
		pairInfo.TotalShares[productId] = pairInfo.TotalShares[productId] - size
		pairInfo.Matches[productId] = append(pairInfo.Matches[productId][:0], pairInfo.Matches[productId][1:]...)
	}
}

func calculateVolumeWeightedAveragePrice(totalSpent, totalSharesBought float64) float64 {
	return totalSpent / totalSharesBought
}
