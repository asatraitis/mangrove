package dal

import (
	"context"
	"testing"
	"time"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type UserTokensDALTestSuite struct {
	suite.Suite

	ctx           context.Context
	DB            *pgxpool.Pool
	userDAL       UserDAL
	userTokensDAL UserTokensDAL

	testUserID uuid.UUID
	testToken  *models.UserToken
}

func TestUserTokensDALTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test suite")
	}
	suite.Run(t, new(UserTokensDALTestSuite))
}

func (suite *UserTokensDALTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	dbpool, err := utils.InitDbPool(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.DB = dbpool

	suite.userTokensDAL = NewUserTokensDAL(suite.ctx, &BaseDAL{
		logger: zerolog.Nop(),
		db:     suite.DB,
	})
	suite.userDAL = NewUserDAL(suite.ctx, &BaseDAL{
		logger: zerolog.Nop(),
		db:     suite.DB,
	})
}

func (suite *UserTokensDALTestSuite) SetupTest() {
	suite.testUserID = uuid.New()
	username := uuid.New().String()

	err := suite.userDAL.Create(nil, &models.User{
		ID:          suite.testUserID,
		Username:    username,
		DisplayName: "Test User Credentials",
		Email:       nil,
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.USER_ROLE_USER,
	})
	suite.NoError(err)
}
func (suite *UserTokensDALTestSuite) TearDownTest() {}

func (suite *UserTokensDALTestSuite) TestCreateGetUserToken_OK() {
	testToken := &models.UserToken{
		ID:      uuid.New(),
		UserID:  suite.testUserID,
		Expires: time.Now().Add(time.Hour * 24),
	}

	err := suite.userTokensDAL.Create(nil, testToken)
	suite.NoError(err)

	createdToken, err := suite.userTokensDAL.Get(testToken.ID)
	suite.NoError(err)
	suite.NotNil(createdToken)

	suite.Equal(testToken.ID.String(), createdToken.ID.String())
	suite.Equal(testToken.UserID.String(), createdToken.UserID.String())

	// TODO: timestamp in PG comes back as UTC vs local time in GO; might need to look into it
	// expected: time.Date(2025, time.January, 18, 9, 45, 18, 930400500, time.Local)
	// actual  : time.Date(2025, time.January, 18, 9, 45, 18, 930400000, time.UTC)
	suite.Equal(testToken.Expires.Format(time.DateTime), createdToken.Expires.Format(time.DateTime))
}
