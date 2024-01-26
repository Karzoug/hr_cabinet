package integration

import (
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (suite *UserTestSuite) TestLogin() {
	t := suite.T()
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	bodyOK := strings.NewReader(`{
		"login": "yapparova@company.com",
		"password": "pa$$word_"
	
	}`)
	respOK, err := retryClient.Post(suite.serverURL()+"/login", "application/json", bodyOK)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respOK.StatusCode)
	assert.Len(t, respOK.Cookies(), 2)

	bodyWrong := strings.NewReader(`{

		"login": "yapparova@company.com",
		"password": "pa$$word_123456"
	
	}`)
	respWrong, err := retryClient.Post(suite.serverURL()+"/login", "application/json", bodyWrong)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, respWrong.StatusCode)
	assert.Len(t, respWrong.Cookies(), 0)
}
