package integration

import (
	"io"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (suite *UserTestSuite) TestCreateVisa() {
	t := suite.T()
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	// login hr
	bodyLogin := strings.NewReader(`{
		"login": "yapparova@company.com",
		"password": "pa$$word_"
	
	}`)
	respLogin, err := retryClient.Post(suite.serverURL()+"/login", "application/json", bodyLogin)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respLogin.StatusCode)

	// get employee visas list before
	respListBefore, err := retryClient.Get(suite.serverURL() + "/users/10/visas")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respListBefore.StatusCode)

	// count visas before
	defer respListBefore.Body.Close()
	bodyListBefore, err := io.ReadAll(respListBefore.Body)
	require.NoError(t, err)
	var visasCountBefore uint
	_, err = jsonparser.ArrayEach(bodyListBefore,
		func(_ []byte, _ jsonparser.ValueType, _ int, _ error) {
			visasCountBefore++
		})
	require.NoError(t, err)

	// add new visa
	bodyCreateVisa := strings.NewReader(`{
		"number": "33592222",
		"issued_state": "Spain",
		"valid_to": "2017-10-22",
		"valid_from": "2017-09-08",
		"type": "C1"
	  }`)
	respCreateVisa, err := retryClient.Post(suite.serverURL()+"/users/10/visas", "application/json", bodyCreateVisa)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, respCreateVisa.StatusCode)
	location := respCreateVisa.Header.Get("Location")
	assert.NotEmpty(t, location)

	// get employee visas list after
	respListAfter, err := retryClient.Get(suite.serverURL() + "/users/10/visas")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respListAfter.StatusCode)

	// count visas after
	defer respListAfter.Body.Close()
	bodyListAfter, err := io.ReadAll(respListAfter.Body)
	require.NoError(t, err)
	var visasCountAfter uint
	_, err = jsonparser.ArrayEach(bodyListAfter,
		func(_ []byte, _ jsonparser.ValueType, _ int, _ error) {
			visasCountAfter++
		})
	require.NoError(t, err)
	assert.Equal(t, visasCountBefore+1, visasCountAfter, "visas count not incremented")
}
