package bll

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/dal/mocks"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ConfigBLLTestSuite struct {
	suite.Suite

	Ctrl      *gomock.Controller
	ctx       context.Context
	logger    zerolog.Logger
	vars      *configs.EnvVariables
	appConfig config.Configs

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
	suite.vars.MangroveSalt = "testsalt"

	suite.configDal = mocks.NewMockConfigDAL(suite.Ctrl)
	suite.dal = mocks.NewMockDAL(suite.Ctrl)
	suite.dal.EXPECT().Config(gomock.Any()).Return(suite.configDal).AnyTimes()

	wauthn, err := webauthn.NewWebAuthN(suite.logger)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.appConfig = config.NewConfig(context.Background(), suite.logger)
	suite.bll = NewBLL(suite.logger, suite.vars, suite.appConfig, wauthn, suite.dal)
}
func (suite *ConfigBLLTestSuite) SetupTest() {
	suite.ctx = context.Background()
}
func (suite *ConfigBLLTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *ConfigBLLTestSuite) TestGetAll_OK() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INIT_SA_CODE: dal.Config{}}, nil).Times(1)

	// run
	configs, err := suite.bll.Config(suite.ctx).GetAll()

	// tests
	suite.NoError(err)
	suite.NotNil(configs)
}

func (suite *ConfigBLLTestSuite) TestGetAll_FAIL() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(nil, errors.New("test error")).Times(1)

	// run
	configs, err := suite.bll.Config(suite.ctx).GetAll()

	// tests
	suite.Error(err)
	suite.ErrorContains(err, "failed to get configs")
	suite.Nil(configs)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCode_OK() {
	// setup
	var val string = "false"
	var attempts string = "0"
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{
		dal.CONFIG_INSTANCE_READY: dal.Config{Key: "instanceReady", Value: &val},
		dal.CONFIG_INIT_SA_CODE:   dal.Config{},
		dal.CONFIG_INIT_ATTEMPTS:  dal.Config{Value: &attempts},
	}, nil).Times(1)
	suite.configDal.EXPECT().Set(dal.CONFIG_INIT_SA_CODE, gomock.Any()).Return(nil)
	suite.configDal.EXPECT().Set(dal.CONFIG_INIT_ATTEMPTS, "0").Return(nil)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.NoError(err)
	suite.NotEqual("", code)
	suite.Equal(6, len(code))

	hashedCode, err := suite.appConfig.GetConfig(dal.CONFIG_INIT_SA_CODE)
	suite.NoError(err)
	suite.NotEmpty(hashedCode)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCode_OK_ReadyInstance() {
	// setup
	val := "true"
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{Key: "instanceReady", Value: &val}}, nil).Times(1)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.NoError(err)
	suite.Equal("", code)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCode_FAIL_noConfig() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(nil, errors.New("test error")).Times(1)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.Error(err)
	suite.Equal("", code)
}

func (suite *ConfigBLLTestSuite) TestInitRegistrationCode_FAIL_emptyConfig() {
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

func (suite *ConfigBLLTestSuite) TestInitRegistrationCode_FAIL_noValue() {
	// setup
	suite.configDal.EXPECT().GetAll().Return(dal.Configs{dal.CONFIG_INSTANCE_READY: dal.Config{Key: "instanceReady"}}, nil).Times(1)

	// run
	code, err := suite.bll.Config(suite.ctx).InitRegistrationCode()

	// test
	suite.Error(err)
	suite.ErrorContains(err, "config instanceReady not set")
	suite.Equal("", code)
}

func (suite *ConfigBLLTestSuite) TestValidateRegistrationCode_OK() {
	// setup
	var testAttempts string = "0"
	var testValue string = "sWCtIzlmuIqI5Q4PWNBtJeXUT/co+a3fZXXVG5Wa8zM=" // 123456
	suite.appConfig.SetAll(dal.Configs{
		dal.CONFIG_INIT_SA_CODE:  dal.Config{Value: &testValue},
		dal.CONFIG_INIT_ATTEMPTS: dal.Config{Value: &testAttempts},
	})

	// run
	err := suite.bll.Config(suite.ctx).ValidateRegistrationCode("123456")

	// test
	suite.NoError(err)
}

func (suite *ConfigBLLTestSuite) TestValidateRegistrationCode_FAIL() {
	// setup
	var testInitAttempts = "0"
	var testValue string = "0arivWpd9loXHq7PRRPZ0svkODQSubIkbW7brExl0mY=" // 175006
	suite.appConfig.SetAll(dal.Configs{
		dal.CONFIG_INIT_SA_CODE:  dal.Config{Value: &testValue},
		dal.CONFIG_INIT_ATTEMPTS: dal.Config{Value: &testInitAttempts},
	})

	// run
	err := suite.bll.Config(suite.ctx).ValidateRegistrationCode("")
	suite.Error(err)
	suite.ErrorContains(err, "no registration code provided")

	// run
	err = suite.bll.Config(suite.ctx).ValidateRegistrationCode("12345")
	suite.Error(err) // expects len of 6
	suite.ErrorContains(err, "no registration code provided")

	// run
	err = suite.bll.Config(suite.ctx).ValidateRegistrationCode("1234557")
	suite.Error(err) // expects len of 6
	suite.ErrorContains(err, "no registration code provided")

	// setup
	suite.configDal.EXPECT().Set(dal.CONFIG_INIT_ATTEMPTS, gomock.Any()).Return(nil)
	// run
	err = suite.bll.Config(suite.ctx).ValidateRegistrationCode("012345")
	fmt.Println(err)
	suite.Error(err) // expects len of 6
	suite.ErrorContains(err, "failed to validate code")
}
