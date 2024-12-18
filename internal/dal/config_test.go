package dal

import (
	"context"
	"testing"

	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ConfigDALTestSuite struct {
	suite.Suite

	ctx       context.Context
	DB        *pgxpool.Pool
	configDAL ConfigDAL
}

func TestConfigDALTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test suite")
	}
	suite.Run(t, new(ConfigDALTestSuite))
}

func (suite *ConfigDALTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	dbpool, err := utils.InitDbPool(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.DB = dbpool

	suite.configDAL = NewConfigDAL(suite.ctx, &BaseDAL{
		logger: zerolog.Nop(),
		db:     suite.DB,
	})
}
func (suite *ConfigDALTestSuite) SetupTest()    {}
func (suite *ConfigDALTestSuite) TearDownTest() {}

func (suite *ConfigDALTestSuite) TestGetAll_OK() {
	conf, err := suite.configDAL.GetAll()
	suite.NoError(err)
	suite.NotNil(conf)
}

func (suite *ConfigDALTestSuite) TestSet_OK() {
	err := suite.configDAL.Set(CONFIG_INSTANCE_READY, "true")
	suite.NoError(err)

	c, err := suite.configDAL.GetAll()
	suite.NoError(err)
	suite.NotNil(c)

	instanceRaedyConf, ok := c[CONFIG_INSTANCE_READY]
	suite.True(ok)
	suite.Equal("true", *instanceRaedyConf.Value)
}
