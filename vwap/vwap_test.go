package vwap

import (
	"VolumePrice/helper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type vWapSuite struct {
	suite.Suite
	*require.Assertions
	pairInfo helper.PairInfo
}

func TestVWapSuite(t *testing.T) {
	suite.Run(t, new(vWapSuite))
}

func (suite *vWapSuite) SetupTest() {
	suite.Assertions = suite.Suite.Require()
	suite.pairInfo.Matches = make(map[string][]helper.Match)
	suite.pairInfo.TotalShares = make(map[string]float64)
	suite.pairInfo.TotalSpent = make(map[string]float64)
	suite.pairInfo.VolumeWeightedAveragePrice = make(map[string]float64)
}

func (suite *vWapSuite) TestVolumeWeightedAveragePriceCalculation() {
	firstCustomer := float64(1000)
	secondCustmoer := float64(2400)
	thirdCustomer := float64(2200)
	totalDollarsSpent := firstCustomer + secondCustmoer + thirdCustomer
	totalSharesPurchased := float64(600)
	vwap := calculateVolumeWeightedAveragePrice(totalDollarsSpent, totalSharesPurchased)
	suite.Equal(9.333333333333334, vwap)
}

func (suite *vWapSuite) TestCalculation() {
	match := helper.Match{ProductId: "BTC-USD", Price: "42931.74", Size: "0.00116661"}
	calculatedInfo := doCalculations(match, suite.pairInfo)
	suite.Equal(calculatedInfo.VolumeWeightedAveragePrice, map[string]float64{"BTC-USD": 42931.74})
}

func (suite *vWapSuite) TestCalculation2() {
	match := helper.Match{ProductId: "BTC-USD", Price: "10", Size: "100"}
	// reference matches for sanity's sake
	//match2 := helper.Match{ProductId: "BTC-USD", Price: "8", Size: "300"}
	//match3 := helper.Match{ProductId: "BTC-USD", Price: "11", Size: "200"}
	suite.pairInfo.TotalSpent["BTC-USD"] = suite.pairInfo.TotalSpent["BTC-USD"] + (float64(300) * float64(8))
	suite.pairInfo.TotalShares["BTC-USD"] = suite.pairInfo.TotalShares["BTC-USD"] + 300
	suite.pairInfo.TotalSpent["BTC-USD"] = suite.pairInfo.TotalSpent["BTC-USD"] + (float64(11) * float64(200))
	suite.pairInfo.TotalShares["BTC-USD"] = suite.pairInfo.TotalShares["BTC-USD"] + 200
	calculatedInfo := doCalculations(match, suite.pairInfo)
	suite.Equal(9.333333333333334, calculatedInfo.VolumeWeightedAveragePrice["BTC-USD"])
}

func (suite *vWapSuite) TestCalculation3() {
	match := helper.Match{ProductId: "ETH-BTC", Price: "0.1", Size: "0.1"}
	calculatedInfo := doCalculations(match, suite.pairInfo)
	suite.Equal(0.10000000000000002, calculatedInfo.VolumeWeightedAveragePrice["ETH-BTC"])
}

func (suite *vWapSuite) TestEmptyCalculation() {
	match := helper.Match{ProductId: "ETH-BTC", Price: "0", Size: "0"}
	calculatedInfo := doCalculations(match, suite.pairInfo)
	suite.NotEmpty(calculatedInfo.VolumeWeightedAveragePrice["ETH-BTC"])
	match2 := helper.Match{ProductId: "ETH-BTC", Price: "0.1", Size: "0.1"}
	calculatedInfo2 := doCalculations(match2, suite.pairInfo)
	suite.Equal(0.10000000000000002, calculatedInfo2.VolumeWeightedAveragePrice["ETH-BTC"])
}

func (suite *vWapSuite) TestRemove() {
	match := helper.Match{ProductId: "BTC-USD", Price: "10", Size: "100"}
	match2 := helper.Match{ProductId: "BTC-USD", Price: "8", Size: "300"}
	match3 := helper.Match{ProductId: "BTC-USD", Price: "11", Size: "200"}
	match4 := helper.Match{ProductId: "BTC-USD", Price: "12", Size: "400"}
	suite.pairInfo.Matches["BTC-USD"] = append(suite.pairInfo.Matches["BTC-USD"], match)
	suite.pairInfo.Matches["BTC-USD"] = append(suite.pairInfo.Matches["BTC-USD"], match2)
	suite.pairInfo.Matches["BTC-USD"] = append(suite.pairInfo.Matches["BTC-USD"], match3)
	suite.pairInfo.TotalSpent["BTC-USD"] = suite.pairInfo.TotalSpent["BTC-USD"] + (float64(300) * float64(8))
	suite.pairInfo.TotalShares["BTC-USD"] = suite.pairInfo.TotalShares["BTC-USD"] + 300
	suite.pairInfo.TotalSpent["BTC-USD"] = suite.pairInfo.TotalSpent["BTC-USD"] + (float64(11) * float64(200))
	suite.pairInfo.TotalShares["BTC-USD"] = suite.pairInfo.TotalShares["BTC-USD"] + 200
	calculatedInfo := doCalculations(match4, suite.pairInfo)
	suite.Equal(10.444444444444445, calculatedInfo.VolumeWeightedAveragePrice["BTC-USD"])
	suite.Equal(4, len(suite.pairInfo.Matches["BTC-USD"]))
	suite.Equal(helper.Match{ProductId: "BTC-USD", Price: "10", Size: "100"}, suite.pairInfo.Matches["BTC-USD"][0])

	removePairCheck(&suite.pairInfo, "BTC-USD", 3)
	suite.Equal(3, len(suite.pairInfo.Matches["BTC-USD"]))
	suite.NotEqual(helper.Match{ProductId: "BTC-USD", Price: "10", Size: "100"}, suite.pairInfo.Matches["BTC-USD"][0])
	suite.Equal(helper.Match{ProductId: "BTC-USD", Price: "8", Size: "300"}, suite.pairInfo.Matches["BTC-USD"][0])
}

func (suite *vWapSuite) TestVWapWorks() {
	matches := make(chan helper.Match, 1)
	match := helper.Match{ProductId: "BTC-USD", Price: "42931.74", Size: "0.00116661"}
	go GetVWap(matches)
	matches <- match
	time.Sleep(time.Second * 1)
	suite.Empty(matches)
}
