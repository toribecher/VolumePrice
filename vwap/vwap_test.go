package vwap

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type vWapSuite struct {
	suite.Suite
	*require.Assertions
}

func TestVWapSuite(t *testing.T) {
	suite.Run(t, new(vWapSuite))
}

func (suite *vWapSuite) SetupTest() {
	suite.Assertions = suite.Suite.Require()
}

func (suite *vWapSuite) TestCalculation() {
	firstCustomer := float64(1000)
	secondCustmoer := float64(2400)
	thirdCustomer := float64(2200)
	totalDollarsSpent := firstCustomer + secondCustmoer + thirdCustomer
	totalSharesPurchased := float64(600)
	vwap := calculateVolumeWeightedAveragePrice(totalDollarsSpent, totalSharesPurchased)
	suite.Equal(9.333333333333334, vwap)
}
