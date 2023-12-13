package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createTargetUrl(city string, count int) string {
	return fmt.Sprintf("/cafe?city=%s&count=%d", city, count)
}

func createRequest(targetUrl string) *http.Request {
	return httptest.NewRequest(http.MethodGet, targetUrl, nil)
}

func TestMainHandlerCorrectRequest(t *testing.T) {
	city := "moscow"
	count := 4
	targetUrl := createTargetUrl(city, count)
	req := createRequest(targetUrl)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.NotEmpty(t, data)
}

func TestMainHandlerUnsupportedCity(t *testing.T) {
	city := "not-moscow"
	count := 4
	targetUrl := createTargetUrl(city, count)
	req := createRequest(targetUrl)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, []byte("wrong city value"), data)
}

func TestMainHandlerWithoutCount(t *testing.T) {
	city := "moscow"

	targetUrl := fmt.Sprintf("/cafe?city=%s", city)
	req := createRequest(targetUrl)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, []byte("count missing"), data)
}

func TestMainHandlerWrongCountValue(t *testing.T) {
	city := "moscow"

	targetUrl := fmt.Sprintf("/cafe?city=%s&count=g", city)
	req := createRequest(targetUrl)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, []byte("wrong count value"), data)
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	city := "moscow"
	count := 5
	totalCount := 4
	targetUrl := createTargetUrl(city, count)
	req := createRequest(targetUrl)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	assert.Equal(t, totalCount, len(strings.Split(string(data), ",")))
}
