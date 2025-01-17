package dal

import (
	"context"
	"testing"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type UserCredentialsDALTestSuite struct {
	suite.Suite

	ctx                context.Context
	DB                 *pgxpool.Pool
	userCredentialsDAL UserCredentialsDAL
	userDAL            UserDAL

	userUUID uuid.UUID
}

func TestUserCredentialsDALTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test suite")
	}
	suite.Run(t, new(UserCredentialsDALTestSuite))
}

func (suite *UserCredentialsDALTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	dbpool, err := utils.InitDbPool(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.DB = dbpool

	suite.userCredentialsDAL = NewUserCredentialsDAL(suite.ctx, &BaseDAL{
		logger: zerolog.Nop(),
		db:     suite.DB,
	})
	suite.userDAL = NewUserDAL(suite.ctx, &BaseDAL{
		logger: zerolog.Nop(),
		db:     suite.DB,
	})
}
func (suite *UserCredentialsDALTestSuite) SetupTest() {
	suite.userUUID = uuid.New()
	username := uuid.New().String()

	err := suite.userDAL.Create(nil, &models.User{
		ID:          suite.userUUID,
		Username:    username,
		DisplayName: "Test User Credentials",
		Email:       nil,
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.USER_ROLE_USER,
	})
	suite.NoError(err)
}
func (suite *UserCredentialsDALTestSuite) TearDownTest() {}

func (suite *UserCredentialsDALTestSuite) TestCreate_OK() {
	userCredential := models.UserCredential{
		ID:                            []byte(uuid.New().String()),
		UserID:                        suite.userUUID,
		PublicKey:                     []byte("test-public-key"),
		AttestationType:               "basic",
		Transport:                     []protocol.AuthenticatorTransport{protocol.USB, protocol.NFC},
		FlagUserPresent:               true,
		FlagVerified:                  true,
		FlagBackupEligible:            true,
		FlagBackupState:               true,
		AuthAaguid:                    []byte("test-aaguid"),
		AuthSignCount:                 uint32(1),
		AuthCloneWarning:              true,
		AuthAttachment:                protocol.CrossPlatform,
		AttestationClientDataJson:     []byte("test-client-data-json"),
		AttestationDataHash:           []byte("test-data-hash"),
		AttestationAuthenticatorData:  []byte("test-authenticator-data"),
		AttestationPublicKeyAlgorithm: int64(1),
		AttestationObject:             []byte("test-attestation-object"),
	}

	err := suite.userCredentialsDAL.Create(nil, &userCredential)
	suite.NoError(err)

	tx, err := suite.DB.BeginTx(suite.ctx, pgx.TxOptions{})
	suite.NoError(err)

	userCredential.ID = []byte(uuid.New().String())
	err = suite.userCredentialsDAL.Create(tx, &userCredential)
	suite.NoError(err)
}

func (suite *UserCredentialsDALTestSuite) TestCreate_FailNoCredential() {
	err := suite.userCredentialsDAL.Create(nil, nil)
	suite.Error(err)
	suite.ErrorContains(err, "nil credential")
}
