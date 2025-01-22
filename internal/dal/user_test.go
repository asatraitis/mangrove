package dal

import (
	"context"
	"testing"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type UserDALTestSuite struct {
	suite.Suite

	ctx         context.Context
	DB          *pgxpool.Pool
	userDAL     UserDAL
	userCredDAL UserCredentialsDAL
}

func TestUserDALTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test suite")
	}
	suite.Run(t, new(UserDALTestSuite))
}

func (suite *UserDALTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	dbpool, err := utils.InitDbPool(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.DB = dbpool

	suite.userDAL = NewUserDAL(suite.ctx, &BaseDAL{
		logger: zerolog.Nop(),
		db:     suite.DB,
	})
	suite.userCredDAL = NewUserCredentialsDAL(suite.ctx, &BaseDAL{
		logger: zerolog.Nop(),
		db:     suite.DB,
	})
}
func (suite *UserDALTestSuite) SetupTest()    {}
func (suite *UserDALTestSuite) TearDownTest() {}

func (suite *UserDALTestSuite) TestCreate_OK() {
	userUUID := uuid.New()
	email := "test@email.com"

	user := &models.User{
		ID:          userUUID,
		Username:    "test",
		DisplayName: "Test User",
		Email:       &email,
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.USER_ROLE_USER,
	}

	err := suite.userDAL.Create(nil, user)
	suite.NoError(err)

	tx, err := suite.DB.BeginTx(suite.ctx, pgx.TxOptions{})
	suite.NoError(err)

	userUUID = uuid.New()
	user.ID = userUUID
	user.Username = "test-with-tx"
	err = suite.userDAL.Create(tx, user)
	suite.NoError(err)
	err = tx.Commit(suite.ctx)
	suite.NoError(err)

	createdUser, err := suite.userDAL.GetByID(userUUID)
	suite.NoError(err)

	suite.Equal(user.ID.String(), createdUser.ID.String())
	suite.Equal(user.Username, createdUser.Username)
	suite.Equal(user.DisplayName, createdUser.DisplayName)
	suite.Equal(user.Status, createdUser.Status)
	suite.Equal(user.Role, createdUser.Role)
}

func (suite *UserDALTestSuite) TestCreate_FailNilUser() {
	err := suite.userDAL.Create(nil, nil)
	suite.Error(err)
	suite.ErrorContains(err, "nil user")
}

func (suite *UserDALTestSuite) TestGetByUsernameWithCredentials_OK() {
	userUUID := uuid.New()
	email := "test@email.com"

	user := &models.User{
		ID:          userUUID,
		Username:    "test" + userUUID.String(),
		DisplayName: "Test User",
		Email:       &email,
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.USER_ROLE_USER,
	}

	err := suite.userDAL.Create(nil, user)
	suite.NoError(err)

	createdUser, err := suite.userDAL.GetByUsernameWithCredentials("test" + userUUID.String())
	suite.NoError(err)

	suite.Equal(user.ID.String(), createdUser.ID.String())
	suite.Equal(user.Username, createdUser.Username)
	suite.Equal(user.DisplayName, createdUser.DisplayName)
	suite.Equal(user.Status, createdUser.Status)
	suite.Equal(user.Role, createdUser.Role)
	suite.Nil(user.Credentials)

	testCredential := getUserCredential(userUUID)
	err = suite.userCredDAL.Create(nil, testCredential)
	suite.NoError(err)

	createdUser, err = suite.userDAL.GetByUsernameWithCredentials("test" + userUUID.String())
	suite.NoError(err)
	suite.NotNil(createdUser.Credentials)
}
