package webauthn

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type WebAuthNTestSuite struct {
	suite.Suite

	logger zerolog.Logger
}

func TestWebAuthNTestSuite(t *testing.T) {
	suite.Run(t, new(WebAuthNTestSuite))
}

func (suite *WebAuthNTestSuite) SetupSuite() {
	suite.logger = zerolog.Nop()
}

func (suite *WebAuthNTestSuite) TestNewWebAuthN() {
	wa, err := NewWebAuthN(suite.logger)
	suite.NoError(err)
	suite.NotNil(wa)
}
