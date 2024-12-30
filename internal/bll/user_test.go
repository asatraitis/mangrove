package bll

import (
	"context"
	"testing"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal/mocks"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserBllTestSuite struct {
	suite.Suite

	Ctrl *gomock.Controller
	ctx  context.Context

	dal *mocks.MockDAL
	bll BLL
}

func TestUserBllTestSuite(t *testing.T) {
	suite.Run(t, new(UserBllTestSuite))
}

func (suite *UserBllTestSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.dal = mocks.NewMockDAL(suite.Ctrl)

	logger := zerolog.Nop()
	vars := configs.NewConf(logger).GetEnvironmentVars()
	vars.MangroveEnv = "testsalt"
	wauthn, err := webauthn.NewWebAuthN(logger)
	if err != nil {
		suite.T().Fatal(err)
	}
	appConfig := config.NewConfig(context.Background(), logger)
	suite.bll = NewBLL(logger, vars, appConfig, wauthn, suite.dal)
}
func (suite *UserBllTestSuite) SetupTest() {
	suite.ctx = context.Background()
}
func (suite *UserBllTestSuite) TearDownTest() {}
func (suite *UserBllTestSuite) TestCreateUserSession_OK() {
	creds, err := suite.bll.User(suite.ctx).CreateUserSession()
	suite.NoError(err)
	suite.NotEmpty(creds)
}

// TODO: decide on how to unit test webauthn flow
func (suite *UserBllTestSuite) TestRegisterSuperAdmin_OK() {

}
