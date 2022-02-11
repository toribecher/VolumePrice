package vwap

import (
	"VolumePrice/helper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
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
	match := helper.Match{ProductId: "BTC-USD", Price: "42931.74", Size: "0.00116661"}
	//match2 := helper.Match{ProductId: "ETH-BTC", Price: "0.07094", Size: "0.00131352"}
	//suite.pairInfo.Matches[]
	calculatedInfo := doCalculations(match, suite.pairInfo)
	suite.Equal(42931.74, calculatedInfo.VolumeWeightedAveragePrice["BTC-USD"])
}

func (suite *vWapSuite) TestCalculation3() {
	match := helper.Match{ProductId: "BTC-USD", Price: "42931.74", Size: "0.00116661"}
	//match2 := helper.Match{ProductId: "ETH-BTC", Price: "0.07094", Size: "0.00131352"}
	calculatedInfo := doCalculations(match, suite.pairInfo)
	suite.Equal(42931.74, calculatedInfo.VolumeWeightedAveragePrice["BTC-USD"])
}

func (suite *vWapSuite) TestRemove() {

}

//data {BTC-USD 42931.74 0.00116661}
//{ETH-USD 3045.93 0.00242001}
//{ETH-BTC 0.07094 0.00131352}
//{ETH-USD 3045.92 0.001}
//{ETH-USD 3045.41 0.00312232}
//{ETH-USD 3045.25 0.00266478}
//{ETH-USD 3045.24 0.44113937}
//{ETH-USD 3045.24 0.0005}
//{ETH-USD 3045.25 0.20214651}
//{ETH-USD 3045.25 0.15990915}
//{ETH-USD 3045.25 0.07408304}
//{ETH-USD 3045.26 0.18811}
//{ETH-USD 3045.26 0.00000808}
//{ETH-USD 3045.42 0.00099192}
//{BTC-USD 42933.99 0.05208915}
//{BTC-USD 42933.99 0.02328879}
//{BTC-USD 42933.04 0.04420668}
//{BTC-USD 42932.38 0.05}
