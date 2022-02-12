package ws

import (
	"VolumePrice/helper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	_ "github.com/stretchr/testify/suite"
	"testing"
)

type webSocketSuite struct {
	suite.Suite
	*require.Assertions
}

func TestWebSocketSuite(t *testing.T) {
	suite.Run(t, new(webSocketSuite))
}

func (suite *webSocketSuite) SetupTest() {
	suite.Assertions = suite.Suite.Require()
}

func (suite *webSocketSuite) TestSubscribe() {
	var websocket WebSocket
	matches := make(chan helper.Match)
	err := websocket.SubscribeAndRead(matches)
	suite.NoError(err)
}

func (suite *webSocketSuite) TestRead() {
	var websocket WebSocket
	matches := make(chan helper.Match)
	websocket.SubscribeAndRead(matches)
	var count int
	var checkArray []helper.Match
	for match := range matches {
		count++
		checkArray = append(checkArray, match)
		if count > 1 {
			break
		}
	}
	suite.NotEmpty(checkArray)
}
