package helper

type Match struct {
	ProductId string `json:"product_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
}

type PairInfo struct {
	VolumeWeightedAveragePrice map[string]float64
	TotalSpent                 map[string]float64
	TotalShares                map[string]float64
	Matches                    map[string][]Match
}
