package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CryptoTestSuite struct {
	suite.Suite

	salt []byte
}

func TestCryptoTestSuite(t *testing.T) {
	suite.Run(t, new(CryptoTestSuite))
}
func (suite *CryptoTestSuite) SetupSuite() {
	suite.salt = []byte("YaaGVyKXtIGl8U4WDamUb8BAKJaiTpz6nPxmND9zNX0=")
}
func (suite *CryptoTestSuite) SetupTest()    {}
func (suite *CryptoTestSuite) TearDownTest() {}

func (suite *CryptoTestSuite) TestEncodeToString() {
	testString := EncodeToString(1)
	suite.NotEqual("", testString)
	suite.Equal(1, len(testString))

	testString = EncodeToString(2)
	suite.NotEqual("", testString)
	suite.Equal(2, len(testString))

	testString = EncodeToString(20)
	suite.NotEqual("", testString)
	suite.Equal(20, len(testString))

	testString = EncodeToString(0)
	suite.Equal("", testString)
	suite.Equal(0, len(testString))

	testString = EncodeToString(-1)
	suite.Equal("", testString)
	suite.Equal(0, len(testString))
}

func (suite *CryptoTestSuite) TestCrypto() {
	c := NewCrypto(1, suite.salt, 64*1024, 4, 32)

	testHash := c.GenerateBase64String([]byte("testHash"))
	suite.NotEqual(32, len(testHash))

	testBytes, err := c.DecodeBase64String(testHash)
	suite.NoError(err)

	err = c.CompareValueToHash("testHash", testBytes)
	suite.NoError(err)
}

func (suite *CryptoTestSuite) TestHmac() {
	c := NewCrypto(1, suite.salt, 64*1024, 4, 32)

	sig := c.GenerateTokenHMAC("test")
	suite.NotEmpty(sig)

	err := c.VerifyToken("test", sig)
	suite.NoError(err)

	err = c.VerifyToken("whatever", sig)
	suite.Error(err)

	err = c.VerifyToken("test", "whatever")
	suite.Error(err)
}
