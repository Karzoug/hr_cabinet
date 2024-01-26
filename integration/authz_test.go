package integration

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (suite *UserTestSuite) TestAuthorization() {
	t := suite.T()

	tests := []struct {
		name           string
		login          string
		password       string
		forbiddenPaths []string
		permittedPaths []string
	}{
		{
			name:           "hr",
			login:          "yapparova@company.com",
			password:       "pa$$word_",
			forbiddenPaths: []string{}, // /accounts
			permittedPaths: []string{"/users", "/users/1"},
		},
		{
			name:           "employee",
			login:          "daryushina@company.com",
			password:       "pa$$word_",
			forbiddenPaths: []string{"/users", "/users/1"}, // /accounts
			permittedPaths: []string{"/users/4"},
		},
		{
			name:           "admin",
			login:          "korepanov@company.com",
			password:       "pa$$word_",
			forbiddenPaths: []string{"/users", "/users/4"},
			permittedPaths: []string{"/users/1"}, // /accounts
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryClient := retryablehttp.NewClient()
			retryClient.RetryMax = 5

			loginBody := strings.NewReader(
				fmt.Sprintf(`{"login": "%s", "password": "%s"}`,
					tt.login, tt.password))
			respLogin, err := retryClient.Post(
				suite.serverURL()+"/login", "application/json",
				loginBody)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, respLogin.StatusCode)
			assert.Len(t, respLogin.Cookies(), 2)

			for _, path := range tt.permittedPaths {
				resp, err := retryClient.Get(suite.serverURL() + path)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}

			for _, path := range tt.forbiddenPaths {
				resp, err := retryClient.Get(suite.serverURL() + path)
				require.NoError(t, err)
				assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			}
		})
	}
}
