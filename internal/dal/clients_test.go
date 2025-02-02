package dal

import (
	"context"
	"testing"
	"time"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ClientsDALTestSuite struct {
	suite.Suite

	ctx context.Context
	DB  *pgxpool.Pool
	dal DAL

	userID uuid.UUID
}

func TestClientsDALTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test suite")
	}
	suite.Run(t, new(ClientsDALTestSuite))
}

func (suite *ClientsDALTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	dbpool, err := utils.InitDbPool(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.DB = dbpool

	dal := NewDAL(zerolog.Nop(), suite.DB)
	suite.dal = dal
}

func (suite *ClientsDALTestSuite) SetupTest() {
	userID := uuid.New()
	email := "test@email.com"

	user := &models.User{
		ID:          userID,
		Username:    "test" + userID.String(),
		DisplayName: "Test User",
		Email:       &email,
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.USER_ROLE_USER,
	}

	err := suite.dal.User(suite.ctx).Create(nil, user)
	suite.NoError(err)

	suite.userID = userID
}
func (suite *ClientsDALTestSuite) TearDownTest() {}

func (suite *ClientsDALTestSuite) TestCreate_OK() {
	ID, err := uuid.NewV7()
	suite.NoError(err)

	expiresAt := time.Now().Add(time.Hour)

	client := models.Client{
		ID:           ID,
		UserID:       suite.userID,
		Name:         "test-client-name",
		Description:  "test-client-description",
		RedirectURI:  "http://localhost:3030",
		PublicKey:    []byte("test-public-key"),
		KeyExpiresAt: expiresAt,
		KeyAlgo:      models.CLIENT_KEY_ALGO_EDDSA,
		Status:       models.CLIENT_STATUS_ACTIVE,
	}

	err = suite.dal.Client(suite.ctx).Create(nil, &client)
	suite.NoError(err)

	tx, err := suite.DB.BeginTx(suite.ctx, pgx.TxOptions{})
	suite.NoError(err)

	ID, err = uuid.NewV7()
	suite.NoError(err)

	client.ID = ID
	err = suite.dal.Client(suite.ctx).Create(tx, &client)
	suite.NoError(err)
}

func (suite *ClientsDALTestSuite) TestCreate_FAIL_NoClient() {
	err := suite.dal.Client(suite.ctx).Create(nil, nil)
	suite.Error(err)
}

func (suite *ClientsDALTestSuite) TestGetAllByUserID_OK() {
	ID, err := uuid.NewV7()
	suite.NoError(err)

	expiresAt := time.Now().Add(time.Hour)

	client := models.Client{
		ID:           ID,
		UserID:       suite.userID,
		Name:         "test-client-name",
		Description:  "test-client-description",
		RedirectURI:  "http://localhost:3030",
		PublicKey:    []byte("test-public-key"),
		KeyExpiresAt: expiresAt,
		KeyAlgo:      models.CLIENT_KEY_ALGO_EDDSA,
		Status:       models.CLIENT_STATUS_ACTIVE,
	}

	err = suite.dal.Client(suite.ctx).Create(nil, &client)
	suite.NoError(err)

	clients, err := suite.dal.Client(suite.ctx).GetAllByUserID(suite.userID)
	suite.NoError(err)
	suite.NotNil(clients)
	suite.NotEqual(0, len(clients))

	var createdClient *models.Client
	for _, created := range clients {
		if ID == created.ID {
			createdClient = created
		}
	}

	suite.NotNil(createdClient)
	suite.Equal(ID, createdClient.ID)
	suite.Equal(suite.userID, createdClient.UserID)
	suite.Equal("test-client-name", createdClient.Name)
	suite.Equal("test-client-description", createdClient.Description)
	suite.Equal("http://localhost:3030", createdClient.RedirectURI)
	suite.Equal(models.ClientStatus("active"), createdClient.Status)
}
