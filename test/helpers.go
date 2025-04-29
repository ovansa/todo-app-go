package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateTestRequest(t *testing.T, router http.Handler, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var requestBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(t, err)
		requestBody = bytes.NewBuffer(jsonData)
	}

	req := httptest.NewRequest(method, path, requestBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), target))
}
