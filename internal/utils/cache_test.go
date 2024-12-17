package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func (suite *CacheTestSuite) TestNewCache() {
	cache := NewCache[string, string]()
	cache.SetValue("testKey", "testValue")
	testValue := cache.GetValue("testKey")

	suite.Equal("testValue", testValue)
}
