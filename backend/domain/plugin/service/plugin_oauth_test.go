package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/coze-dev/coze-studio/backend/domain/plugin/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type pluginOAuthSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	ctx  context.Context

	mockOauthRepo *mock_plugin_oauth.MockOAuthRepository
}

func TestPluginOAuthSuite(t *testing.T) {
	suite.Run(t, &pluginOAuthSuite{})
}

func (s *pluginOAuthSuite) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())
	s.mockOauthRepo = mock_plugin_oauth.NewMockOAuthRepository(s.ctrl)
}

func (s *pluginOAuthSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *pluginOAuthSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *pluginOAuthSuite) TearDownTest() {

}

func (s *pluginOAuthSuite) TestRefreshTokenFailedHandler() {
	mockRecordID := int64(123)
	mockErr := fmt.Errorf("mock error")
	mockSVC := &pluginServiceImpl{
		oauthRepo: s.mockOauthRepo,
	}

	mockSVC.refreshTokenFailedHandler(s.ctx, mockRecordID, mockErr)
	failedTimes, ok := failedCache.Load(mockRecordID)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), 1, failedTimes.(int))

	for i := 2; i < 5; i++ {
		mockSVC.refreshTokenFailedHandler(s.ctx, mockRecordID, mockErr)
		failedTimes, ok = failedCache.Load(mockRecordID)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), i, failedTimes.(int))
	}

	s.mockOauthRepo.EXPECT().BatchDeleteAuthorizationCodeByIDs(gomock.Any(), gomock.Any()).
		Return(nil).Times(1)

	mockSVC.refreshTokenFailedHandler(s.ctx, mockRecordID, mockErr)
	_, ok = failedCache.Load(mockRecordID)
	assert.False(s.T(), ok)
}
