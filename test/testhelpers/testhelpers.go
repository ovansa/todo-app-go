package testhelpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func CreateTestRequest(
	t *testing.T,
	router *gin.Engine,
	method string,
	url string,
	body interface{},
	token string,
) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonBody, err := json.Marshal(body)
		assert.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, dest interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), dest)
	assert.NoError(t, err)
}
