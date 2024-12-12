package bll

import (
	"context"
	"errors"
	"testing"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/dal/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ConfigBLLTestSuite struct {
	suite.Suite

	Ctrl   *gomock.Controller
	ctx    context.Context
	logger zerolog.Logger
	vars   *configs.EnvVariables

	dal       *mocks.MockDAL
	configDal *mocks.MockConfigDAL
	bll       BLL
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigBLLTestSuite))
}
func (suite *ConfigBLLTestSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.logger = zerolog.Nop()
	suite.vars = configs.NewConf(suite.logger).GetEnvironmentVars()

	suite.configDal = mocks.NewMockConfigDAL(suite.Ctrl)
	suite.dal = mocks.NewMockDAL(suite.Ctrl)
	suite.dal.EXPECT().Config(gomock.Any()).Return(suite.configDal).AnyTimes()

	suite.bll = NewBLL(suite.logger, suite.vars, suite.dal)
}
func (suite *ConfigBLLTestSuite) SetupTest() {
	suite.ctx = context.Background()

}
func (suite *ConfigBLLTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *ConfigBLLTestSuite) TestGetAllOK() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INIT_SA_CODE: dal.Config{}}, nil).Times(1)

	// run
	configs, err := suite.bll.Config(suite.ctx).GetAll()

	// tests
	suite.NoError(err)
	suite.NotNil(configs)
}

func (suite *ConfigBLLTestSuite) TestGetAllFAIL() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(nil, errors.New("test error")).Times(1)

	// run
	configs, err := suite.bll.Config(suite.ctx).GetAll()

	// tests
	suite.Error(err)
	suite.ErrorContains(err, "failed to get configs")
	suite.Nil(configs)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCodeOK() {
	// setup
	val := "false"
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{Key: "instanceReady", Value: &val}}, nil).Times(1)
	suite.configDal.EXPECT().Set(dal.CONFIG_INIT_SA_CODE, gomock.Any()).Return(nil)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.NoError(err)
	suite.NotEqual("", code)
	suite.Equal(6, len(code))
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCodeOKReadyInstance() {
	// setup
	val := "true"
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{Key: "instanceReady", Value: &val}}, nil).Times(1)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.NoError(err)
	suite.Equal("", code)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCodeFAILnoConfig() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(nil, errors.New("test error")).Times(1)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.Error(err)
	suite.Equal("", code)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCodeFAILEmptyConfig() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(nil, nil).Times(1)
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{}, nil).Times(1)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config instanceReady was not retrieved")
	suite.Equal("", code)

	// run 2
	code, err = suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config instanceReady not set")
	suite.Equal("", code)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCodeFAILnoValue() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{Key: "instanceReady"}}, nil).Times(1)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config instanceReady not set")
	suite.Equal("", code)
}
