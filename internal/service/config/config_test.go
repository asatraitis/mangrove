package config

import (
	"context"
	"testing"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ConfigServiceTestSuite struct {
	suite.Suite

	ctx    context.Context
	logger zerolog.Logger

	configs Configs
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigServiceTestSuite))
}
func (suite *ConfigServiceTestSuite) SetupSuite() {
	suite.logger = zerolog.Nop()
	suite.ctx = context.Background()
}

func (suite *ConfigServiceTestSuite) SetupTest() {
	suite.configs = NewConfig(suite.ctx, suite.logger)
}
func (suite *ConfigServiceTestSuite) TearDownTest() {
}

func (suite *ConfigServiceTestSuite) TestGetAll_OK() {
	// setup
	var testVal string = "testValue"
	suite.configs.SetAll(dal.Configs{
		dal.CONFIG_INIT_SA_CODE: dal.Config{Key: "initSACode", Value: &testVal},
	})

	// run
	configs := suite.configs.GetAll()

	// test
	suite.NotNil(configs)
	suite.Equal("testValue", *configs[dal.CONFIG_INIT_SA_CODE].Value)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_OK() {
	// setup
	var value string = "123456"
	suite.configs.SetAll(dal.Configs{dal.CONFIG_INIT_SA_CODE: dal.Config{Value: &value}})

	// run
	configVal, err := suite.configs.GetConfig("initSACode")

	// test
	suite.NoError(err)
	suite.Equal("123456", configVal)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_FAIL_nilConfigs() {
	// run
	configVal, err := suite.configs.GetConfig("test")

	// test
	suite.Error(err)
	suite.Equal("", configVal)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_FAIL_configNotExists() {
	// setup
	suite.configs.SetAll(dal.Configs{})

	// run
	configVal, err := suite.configs.GetConfig("test")

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config not found")
	suite.Equal("", configVal)
}

func (suite *ConfigServiceTestSuite) TestGetConfig_FAIL_configValueNotExists() {
	// setup
	suite.configs.SetAll(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{}})

	// run
	configVal, err := suite.configs.GetConfig("instanceReady")

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config not set (nil)")
	suite.Equal("", configVal)
}

func (suite *ConfigServiceTestSuite) TestUpdateConfig_OK() {
	// setup
	suite.configs.SetAll(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{}})

	// run
	err := suite.configs.UpdateConfig(dal.CONFIG_INSTANCE_READY, "testValue")
	suite.NoError(err)

	newVal, err := suite.configs.GetConfig(dal.CONFIG_INSTANCE_READY)
	suite.NoError(err)
	suite.Equal("testValue", newVal)
}

func (suite *ConfigServiceTestSuite) TestUpdateConfig_FAIL() {
	// setup
	suite.configs.SetAll(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{}})

	// run
	err := suite.configs.UpdateConfig("", "testValue")
	suite.Error(err)
	suite.ErrorContains(err, "empty key or value")

	err = suite.configs.UpdateConfig(dal.CONFIG_INSTANCE_READY, "")
	suite.Error(err)
	suite.ErrorContains(err, "empty key or value")

	err = suite.configs.UpdateConfig("whatever", "testValue")
	suite.Error(err)
	suite.ErrorContains(err, "config does not exist")
}
