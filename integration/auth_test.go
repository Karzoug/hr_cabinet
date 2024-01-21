package integration

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/buger/jsonparser"
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

func (suite *UserTestSuite) TestInitChangePassword() {
	t := suite.T()
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	// init change password
	body := strings.NewReader(`{
		"login": "yapparova@company.com"
	}`)
	respInit, err := retryClient.Post(suite.serverURL()+"/login/init-change-password", "application/json", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respInit.StatusCode)

	// try to find key in mail
	textMail, err := suite.findLastMailTo("yapparova@company.com")
	require.NoError(t, err)

	splitted := strings.SplitAfter(textMail, "password-reset?key=")
	require.Len(t, splitted, 2)
	key := splitted[1][0:36]

	// verify key
	respCheckKey, err := retryClient.Get(suite.serverURL() + "/login/change-password?key=" + key)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCheckKey.StatusCode)
}

func (suite *UserTestSuite) TestChangePassword() {
	t := suite.T()
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	const email = "zelner@company.com"

	// init change password
	bodyInit := strings.NewReader(fmt.Sprintf(`{
		"login": "%s"
	}`, email))
	respInit, err := retryClient.Post(suite.serverURL()+"/login/init-change-password", "application/json", bodyInit)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respInit.StatusCode)

	// try to find key in mail
	textMail, err := suite.findLastMailTo(email)
	require.NoError(t, err)

	splitted := strings.SplitAfter(textMail, "password-reset?key=")
	require.Len(t, splitted, 2)
	key := splitted[1][0:36]

	// change password
	const password = "abcdefghijk"
	bodyChange := strings.NewReader(fmt.Sprintf(`{
		"key": "%s",
		"password": "%s"
	  }`, key, password))
	respChange, err := retryClient.Post(suite.serverURL()+"/login/change-password", "application/json", bodyChange)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respChange.StatusCode)

	// try to login with new password
	bodyOK := strings.NewReader(fmt.Sprintf(`{
		"login": "%s",
		"password": "%s"
	
	}`, email, password))
	respOK, err := retryClient.Post(suite.serverURL()+"/login", "application/json", bodyOK)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respOK.StatusCode)
	assert.Len(t, respOK.Cookies(), 2)
}

func (suite *UserTestSuite) findLastMailTo(to string) (string, error) {
	getList := func() ([]byte, error) {
		resp, err := http.Get(fmt.Sprintf("%s/search?kind=to&query=%s", suite.containers.mailhog.ApiURI(), to))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("wrong status code from mailhog: %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	// waiting for the mail
	var (
		bodyList []byte
		err      error
		count    int64
	)
	for attempt := 0; attempt < 20; attempt++ {
		bodyList, err = getList()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		count, err = jsonparser.GetInt(bodyList, "count")
		if err != nil || count == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
	}
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", fmt.Errorf("not received a single mail from mailhog to %s", to)
	}

	textMail, err := jsonparser.GetString(bodyList, "items", "[0]", "Content", "Body")
	if err != nil {
		return "", fmt.Errorf("wrong response format from mailhog: %w", err)
	}

	return textMail, nil
}
