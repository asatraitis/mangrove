package bll

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal/mocks"
	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/handler/types"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ClientBllTestSuite struct {
	suite.Suite

	Ctrl *gomock.Controller
	ctx  context.Context

	dal        *mocks.MockDAL
	clientsDal *mocks.MockClientsDAL
	bll        BLL
}

func TestClientBllTestSuite(t *testing.T) {
	suite.Run(t, new(ClientBllTestSuite))
}

func (suite *ClientBllTestSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.clientsDal = mocks.NewMockClientsDAL(suite.Ctrl)
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
func (suite *ClientBllTestSuite) SetupTest() {
	suite.ctx = context.Background()
}
func (suite *ClientBllTestSuite) TearDownTest() {}

func (suite *ClientBllTestSuite) TestGetUserClients_OK() {
	suite.ctx = context.WithValue(suite.ctx, types.REQ_CTX_KEY_USER_ID, "0bdd05ec-8008-4869-b6ec-6d812ce95507")
	now := time.Now()
	suite.dal.EXPECT().Client(gomock.Any()).Times(1).Return(suite.clientsDal)
	suite.clientsDal.EXPECT().GetAllByUserID(gomock.Any()).Times(1).Return([]*models.Client{
		{
			ID:           uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95508"),
			UserID:       uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95507"),
			Name:         "test-client-name-0",
			Description:  "test-client-description-0",
			Type:         "app-0",
			RedirectURI:  "http://localhost:3030",
			PublicKey:    []byte("test-public-key-0"),
			KeyExpiresAt: now,
			KeyAlgo:      models.CLIENT_KEY_ALGO_EDDSA,
			Status:       models.CLIENT_STATUS_ACTIVE,
		},
		{
			ID:           uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95509"),
			UserID:       uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95507"),
			Name:         "test-client-name-1",
			Description:  "test-client-description-1",
			Type:         "app-1",
			RedirectURI:  "http://localhost:3031",
			PublicKey:    []byte("test-public-key-1"),
			KeyExpiresAt: now,
			KeyAlgo:      models.CLIENT_KEY_ALGO_EDDSA,
			Status:       models.CLIENT_STATUS_ACTIVE,
		},
	}, nil)

	userClients, err := suite.bll.Client(suite.ctx).GetUserClients()
	suite.NoError(err)
	suite.Len(userClients, 2)

	suite.Equal("0bdd05ec-8008-4869-b6ec-6d812ce95508", userClients[0].ID)
	suite.Equal("0bdd05ec-8008-4869-b6ec-6d812ce95509", userClients[1].ID)
	suite.Equal("0bdd05ec-8008-4869-b6ec-6d812ce95507", userClients[0].UserID)
	suite.Equal("0bdd05ec-8008-4869-b6ec-6d812ce95507", userClients[1].UserID)
	suite.Equal("0bdd05ec-8008-4869-b6ec-6d812ce95507", userClients[0].UserID)
	suite.Equal("test-client-name-0", userClients[0].Name)
	suite.Equal("test-client-name-1", userClients[1].Name)
	suite.Equal("test-client-description-0", userClients[0].Description)
	suite.Equal("test-client-description-1", userClients[1].Description)
	suite.Equal("app-0", userClients[0].Type)
	suite.Equal("app-1", userClients[1].Type)
	suite.Equal("http://localhost:3030", userClients[0].RedirectURI)
	suite.Equal("http://localhost:3031", userClients[1].RedirectURI)
	suite.Equal(dto.UserClientStatus("active"), userClients[0].Status)
	suite.Equal(dto.UserClientStatus("active"), userClients[1].Status)
}

func (suite *ClientBllTestSuite) TestGetUserClients_FAIL_NoUserIDInCtx() {
	_, err := suite.bll.Client(suite.ctx).GetUserClients()
	suite.Error(err)
	suite.ErrorContains(err, "failed type assertion")
}
func (suite *ClientBllTestSuite) TestGetUserClients_FAIL_DalErr() {
	suite.ctx = context.WithValue(suite.ctx, types.REQ_CTX_KEY_USER_ID, "0bdd05ec-8008-4869-b6ec-6d812ce95507")

	suite.dal.EXPECT().Client(gomock.Any()).Times(1).Return(suite.clientsDal)
	suite.clientsDal.EXPECT().GetAllByUserID(gomock.Any()).Times(1).Return(nil, errors.New("test"))

	_, err := suite.bll.Client(suite.ctx).GetUserClients()
	suite.Error(err)
	suite.ErrorContains(err, "test")
}
func (suite *ClientBllTestSuite) TestGetUserClients_FAIL_TypeconvErr() {
	suite.ctx = context.WithValue(suite.ctx, types.REQ_CTX_KEY_USER_ID, "0bdd05ec-8008-4869-b6ec-6d812ce95507")

	suite.dal.EXPECT().Client(gomock.Any()).Times(1).Return(suite.clientsDal)
	suite.clientsDal.EXPECT().GetAllByUserID(gomock.Any()).Times(1).Return(nil, nil)

	_, err := suite.bll.Client(suite.ctx).GetUserClients()
	suite.Error(err)
	suite.ErrorContains(err, "client is nil")
}

func (suite *ClientBllTestSuite) TestCreate_OK() {
	suite.ctx = context.WithValue(suite.ctx, types.REQ_CTX_KEY_USER_ID, "0bdd05ec-8008-4869-b6ec-6d812ce95507")
	clientReq := dto.CreateClientRequest{
		Name:        "test-name",
		Description: "test-desc",
		Type:        "test-type",
		RedirectURI: "http://test.com",
		Status:      dto.UserClientStatus("active"),
		PublicKey:   []byte("pub_key"),
		KeyAlgo:     dto.UserClientKeyAlgo("EdDSA"),
	}
	suite.dal.EXPECT().Client(gomock.Any()).Times(1).Return(suite.clientsDal)
	suite.clientsDal.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	clientRes, err := suite.bll.Client(suite.ctx).Create(clientReq)
	suite.NoError(err)
	suite.NotNil(clientRes)
}

func (suite *ClientBllTestSuite) TestCreate_FAIL_BadPayload() {
	clientRes, err := suite.bll.Client(suite.ctx).Create(dto.CreateClientRequest{})
	suite.Error(err)
	suite.ErrorContains(err, "missing")
	suite.Nil(clientRes)
}

func (suite *ClientBllTestSuite) TestCreate_FAIL_BadUserIDinCtx() {
	suite.ctx = context.WithValue(suite.ctx, types.REQ_CTX_KEY_USER_ID, uint32(1))
	clientReq := dto.CreateClientRequest{
		Name:        "test-name",
		Description: "test-desc",
		Type:        "test-type",
		RedirectURI: "http://test.com",
		Status:      dto.UserClientStatus("active"),
		PublicKey:   []byte("pub_key"),
		KeyAlgo:     dto.UserClientKeyAlgo("EdDSA"),
	}

	clientRes, err := suite.bll.Client(suite.ctx).Create(clientReq)
	suite.Error(err)
	suite.ErrorContains(err, "failed type assertion")
	suite.Nil(clientRes)
}

func (suite *ClientBllTestSuite) TestCreate_FAIL_DalErr() {
	suite.ctx = context.WithValue(suite.ctx, types.REQ_CTX_KEY_USER_ID, "0bdd05ec-8008-4869-b6ec-6d812ce95507")
	clientReq := dto.CreateClientRequest{
		Name:        "test-name",
		Description: "test-desc",
		Type:        "test-type",
		RedirectURI: "http://test.com",
		Status:      dto.UserClientStatus("active"),
		PublicKey:   []byte("pub_key"),
		KeyAlgo:     dto.UserClientKeyAlgo("EdDSA"),
	}
	suite.dal.EXPECT().Client(gomock.Any()).Times(1).Return(suite.clientsDal)
	suite.clientsDal.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("test failed"))

	clientRes, err := suite.bll.Client(suite.ctx).Create(clientReq)
	suite.Error(err)
	suite.ErrorContains(err, "failed")
	suite.Nil(clientRes)
}

func (suite *ClientBllTestSuite) TestValidateCreateReq() {
	clientReq := dto.CreateClientRequest{
		Name:        "test-name",
		Type:        "test-type",
		RedirectURI: "http://test.com",
		Status:      dto.UserClientStatus("active"),
		PublicKey:   []byte("pub_key"),
		KeyAlgo:     dto.UserClientKeyAlgo("EdDSA"),
	}

	err := validateCreateReq(clientReq)
	suite.NoError(err)

	err = validateCreateReq(dto.CreateClientRequest{})
	suite.Error(err)
	suite.ErrorContains(err, "missing name")
	suite.ErrorContains(err, "missing type")
	suite.ErrorContains(err, "missing redirectUri")
	suite.ErrorContains(err, "missing or wrong status")
	suite.ErrorContains(err, "missing publicKey")
	suite.ErrorContains(err, "missing or wrong keyAlgo")

	clientReq.Name = " "
	clientReq.Status = dto.UserClientStatus(" ")
	err = validateCreateReq(clientReq)
	suite.Error(err)
	suite.ErrorContains(err, "missing name")
	suite.ErrorContains(err, "missing or wrong status")
}
