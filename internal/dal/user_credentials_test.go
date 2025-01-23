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
	userCredential := getUserCredential(suite.userUUID)

	err := suite.userCredentialsDAL.Create(nil, userCredential)
	suite.NoError(err)

	tx, err := suite.DB.BeginTx(suite.ctx, pgx.TxOptions{})
	suite.NoError(err)

	userCredential.ID = []byte(uuid.New().String())
	err = suite.userCredentialsDAL.Create(tx, userCredential)
	suite.NoError(err)
}

func (suite *UserCredentialsDALTestSuite) TestCreate_FailNoCredential() {
	err := suite.userCredentialsDAL.Create(nil, nil)
	suite.Error(err)
	suite.ErrorContains(err, "nil credential")
}

func (suite *UserCredentialsDALTestSuite) TestGetByUserID_OK() {
	userCredential := getUserCredential(suite.userUUID)
	err := suite.userCredentialsDAL.Create(nil, userCredential)
	suite.NoError(err)

	creds, err := suite.userCredentialsDAL.GetByUserID(suite.userUUID)
	suite.NoError(err)
	suite.Len(creds, 1)

	suite.Equal([]byte(suite.userUUID.String()), creds[0].ID)
	suite.Equal(suite.userUUID.String(), creds[0].UserID.String())
	suite.Equal([]byte("test-public-key"), creds[0].PublicKey)
	suite.Equal("basic", creds[0].AttestationType)

	suite.Len(creds[0].Transport, 2)
	suite.Equal(protocol.AuthenticatorTransport("usb"), creds[0].Transport[0])
	suite.Equal(protocol.AuthenticatorTransport("nfc"), creds[0].Transport[1])

	suite.Equal(true, creds[0].FlagUserPresent)
	suite.Equal(true, creds[0].FlagVerified)
	suite.Equal(true, creds[0].FlagBackupEligible)
	suite.Equal(true, creds[0].FlagBackupState)

	suite.Equal([]byte("test-aaguid"), creds[0].AuthAaguid)
	suite.Equal(uint32(1), creds[0].AuthSignCount)
	suite.Equal(true, creds[0].AuthCloneWarning)
	suite.Equal(protocol.CrossPlatform, creds[0].AuthAttachment)

	suite.Equal([]byte("test-client-data-json"), creds[0].AttestationClientDataJson)
	suite.Equal([]byte("test-data-hash"), creds[0].AttestationDataHash)
	suite.Equal([]byte("test-authenticator-data"), creds[0].AttestationAuthenticatorData)
	suite.Equal(int64(1), creds[0].AttestationPublicKeyAlgorithm)
	suite.Equal([]byte("test-attestation-object"), creds[0].AttestationObject)
}

func (suite *UserCredentialsDALTestSuite) TestGetByUserID_BadID() {
	creds, err := suite.userCredentialsDAL.GetByUserID(uuid.New())
	suite.NoError(err)
	suite.Len(creds, 0)
}

func (suite *UserCredentialsDALTestSuite) TestUpdateSignCount_OK() {
	userCredential := getUserCredential(suite.userUUID)
	err := suite.userCredentialsDAL.Create(nil, userCredential)
	suite.NoError(err)

	creds, err := suite.userCredentialsDAL.GetByUserID(suite.userUUID)
	suite.NoError(err)
	suite.Len(creds, 1)
	suite.Equal(uint32(1), creds[0].AuthSignCount)

	err = suite.userCredentialsDAL.UpdateSignCount(nil, userCredential.ID, 2)
	suite.NoError(err)

	creds, err = suite.userCredentialsDAL.GetByUserID(suite.userUUID)
	suite.NoError(err)
	suite.Len(creds, 1)
	suite.Equal(uint32(2), creds[0].AuthSignCount)
}
