package handler

import (
	"net/http"
	"testing"

	"github.com/asatraitis/mangrove/internal/bll/mocks"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type InitHandlerTestSuite struct {
	suite.Suite
	Ctrl *gomock.Controller

	logger zerolog.Logger
	bll    *mocks.MockBLL
	mux    *http.ServeMux

	initHandler InitHandler
}

func TestInitHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(InitHandlerTestSuite))
}

func (suite *InitHandlerTestSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.logger = zerolog.Nop()
	suite.bll = mocks.NewMockBLL(suite.Ctrl)

	suite.initHandler = NewInitHandler(suite.logger, suite.bll, http.NewServeMux())
}
func (suite *InitHandlerTestSuite) SetupTest() {}
func (suite *InitHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *InitHandlerTestSuite) TestHome_OK() {
	w := utils.NewMockResponseWriter()
	r, err := http.NewRequest("GET", "/", http.NoBody)

	suite.NoError(err)

	suite.initHandler.home(w, r)

	suite.Equal(200, w.Code)
}
