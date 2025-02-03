package bll

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal/mocks"
	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	wa "github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserBllTestSuite struct {
	suite.Suite

	Ctrl *gomock.Controller
	ctx  context.Context

	dal          *mocks.MockDAL
	userDal      *mocks.MockUserDAL
	userTokenDal *mocks.MockUserTokensDAL
	bll          BLL
}

func TestUserBllTestSuite(t *testing.T) {
	suite.Run(t, new(UserBllTestSuite))
}

func (suite *UserBllTestSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.userDal = mocks.NewMockUserDAL(suite.Ctrl)
	suite.userTokenDal = mocks.NewMockUserTokensDAL(suite.Ctrl)
	suite.dal = mocks.NewMockDAL(suite.Ctrl)

	logger := zerolog.Nop()
	vars := configs.NewConf(logger).GetEnvironmentVars()
	vars.MangroveEnv = "testsalt"
	wauthn, err := webauthn.NewWebAuthN(&wa.Config{
		RPDisplayName: "Mangrove",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3030", "http://localhost:3000"},
	}, logger)
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

func (suite *UserBllTestSuite) TestGetUserByID_OK() {
	testUser := models.User{
		ID:          uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95507"),
		Username:    "test-username",
		DisplayName: "test-displayName",
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.USER_ROLE_USER,
	}
	suite.dal.EXPECT().User(gomock.Any()).Times(1).Return(suite.userDal)
	suite.userDal.EXPECT().GetByID(testUser.ID).Times(1).Return(&testUser, nil)

	user, err := suite.bll.User(suite.ctx).GetUserByID(testUser.ID)
	suite.NoError(err)
	suite.Equal("0bdd05ec-8008-4869-b6ec-6d812ce95507", user.ID.String())
	suite.Equal("test-username", user.Username)
	suite.Equal("test-displayName", user.DisplayName)
	suite.Equal(models.UserStatus("active"), user.Status)
	suite.Equal(models.UserRole("user"), user.Role)
}

func (suite *UserBllTestSuite) TestGetUserByID_FAIL() {

	suite.dal.EXPECT().User(gomock.Any()).Times(1).Return(suite.userDal)
	suite.userDal.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("test error"))

	user, err := suite.bll.User(suite.ctx).GetUserByID(uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95507"))
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UserBllTestSuite) TestValidateTokenAndGetUser_OK() {
	suite.dal.EXPECT().UserTokens(gomock.Any()).Times(1).Return(suite.userTokenDal)
	suite.userTokenDal.EXPECT().GetByIdWithUser(uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117a")).Times(1).Return(&models.UserToken{
		ID:      uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117a"),
		UserID:  uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117b"),
		Expires: time.Now().Add(time.Hour * 24),
		User: &models.User{
			ID:          uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117b"),
			Username:    "test-user",
			DisplayName: "test display name",
			Status:      models.UserStatus("active"),
			Role:        models.UserRole("user"),
		},
	}, nil)

	user, err := suite.bll.User(suite.ctx).ValidateTokenAndGetUser(uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117a"))
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal("561fe1de-21dd-45a7-91f6-b2d831fe117b", user.ID.String())
	suite.Equal("test-user", user.Username)
	suite.Equal("test display name", user.DisplayName)
	suite.Equal(models.UserStatus("active"), user.Status)
	suite.Equal(models.UserRole("user"), user.Role)
}

func (suite *UserBllTestSuite) TestValidateTokenAndGetUser_FAIL_ExpiredToken() {
	suite.dal.EXPECT().UserTokens(gomock.Any()).Times(1).Return(suite.userTokenDal)
	suite.userTokenDal.EXPECT().GetByIdWithUser(uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117a")).Times(1).Return(&models.UserToken{
		ID:      uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117a"),
		UserID:  uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117b"),
		Expires: time.Now().Add(time.Hour * -1),
		User: &models.User{
			ID:          uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117b"),
			Username:    "test-user",
			DisplayName: "test display name",
			Status:      models.UserStatus("active"),
			Role:        models.UserRole("user"),
		},
	}, nil)

	user, err := suite.bll.User(suite.ctx).ValidateTokenAndGetUser(uuid.MustParse("561fe1de-21dd-45a7-91f6-b2d831fe117a"))
	suite.Error(err)
	suite.Nil(user)
}

// TODO: decide on how to unit test webauthn flow
func (suite *UserBllTestSuite) TestRegisterSuperAdmin_OK() {

}
func (suite *UserBllTestSuite) TestInitLogin_OK() {}
