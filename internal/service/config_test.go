package service

import (
	"context"
	"errors"
	"testing"

	bllMocks "github.com/asatraitis/mangrove/internal/bll/mocks"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ConfigServiceTestSuite struct {
	suite.Suite
	Ctrl *gomock.Controller

	ctx       context.Context
	logger    zerolog.Logger
	bll       *bllMocks.MockBLL
	configBll *bllMocks.MockConfigBLL

	configs Configs
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigServiceTestSuite))
}
func (suite *ConfigServiceTestSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.logger = zerolog.Nop()
	suite.ctx = context.Background()
	suite.configBll = bllMocks.NewMockConfigBLL(suite.Ctrl)
	suite.bll = bllMocks.NewMockBLL(suite.Ctrl)
	suite.bll.EXPECT().Config(gomock.Any()).Return(suite.configBll).AnyTimes()
}

func (suite *ConfigServiceTestSuite) SetupTest() {
	suite.configs = NewConfig(suite.ctx, suite.logger, suite.bll)
}
func (suite *ConfigServiceTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *ConfigServiceTestSuite) TestGetAll_OK() {
	// setup
	suite.configBll.EXPECT().GetAll().Return(dal.Configs{}, nil).Times(1)

	// run
	configs, err := suite.configs.GetAll()

	// test
	suite.NoError(err)
	suite.NotNil(configs)

	// run 2 (configs from first run should be cached; should not call BLL again)
	configs, err = suite.configs.GetAll()

	// test
	suite.NoError(err)
	suite.NotNil(configs)
}

func (suite *ConfigServiceTestSuite) TestGetAll_FAIL_bllError() {
	// setup
	suite.configBll.EXPECT().GetAll().Return(nil, errors.New("test error")).Times(1)

	// run
	configs, err := suite.configs.GetAll()

	// test
	suite.Nil(configs)
	suite.Error(err)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_OK() {
	// setup
	var value string = "123456"
	suite.configBll.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INIT_SA_CODE: dal.Config{Value: &value}}, nil).Times(1)

	// run
	configVal, err := suite.configs.GetConfig("initSACode")

	// test
	suite.NoError(err)
	suite.Equal("123456", configVal)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_FAIL_nilConfigs() {
	// setup
	suite.configBll.EXPECT().GetAll().Return(nil, nil).Times(1)

	// run
	configVal, err := suite.configs.GetConfig("test")

	// test
	suite.Error(err)
	suite.Equal("", configVal)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_FAIL_bllError() {
	// setup
	suite.configBll.EXPECT().GetAll().Return(nil, errors.New("test err")).Times(1)

	// run
	configVal, err := suite.configs.GetConfig("test")

	// test
	suite.Error(err)
	suite.Equal("", configVal)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_FAIL_configNotExists() {
	// setup
	suite.configBll.EXPECT().GetAll().Return(dal.Configs{}, nil).Times(1)

	// run
	configVal, err := suite.configs.GetConfig("test")

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config not found")
	suite.Equal("", configVal)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_FAIL_configValueNotExists() {
	// setup
	suite.configBll.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{}}, nil).Times(1)

	// run
	configVal, err := suite.configs.GetConfig("instanceReady")

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config not set (nil)")
	suite.Equal("", configVal)
}

func (suite *ConfigServiceTestSuite) TestReload_OK() {
	// setup
	var initialVal string = "testValue"
	suite.configBll.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{Value: &initialVal}}, nil).Times(1)

	// run
	configVal, err := suite.configs.GetConfig("instanceReady")

	// test
	suite.NoError(err)
	suite.Equal("testValue", configVal)

	// run 2
	configVal, err = suite.configs.GetConfig("instanceReady")

	// test 2
	suite.NoError(err)
	suite.Equal("testValue", configVal)

	// setup new value for reload
	var reloadVal string = "testValue#2"
	suite.configBll.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{Value: &reloadVal}}, nil).Times(1)

	// run 3
	err = suite.configs.Reload()

	// test 3
	suite.NoError(err)

	// run 4
	configValReloaded, err := suite.configs.GetConfig("instanceReady")

	// test 4
	suite.NoError(err)
	suite.NotEqual("", configValReloaded)
	suite.NotEqual(configVal, configValReloaded)
	suite.Equal("testValue#2", configValReloaded)

}
