package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	//Arrange
	testTable := []struct {
		name             string
		respMsg          string
		statusCode       int
		expectedRespBody string
	}{
		{
			name:             "test ok",
			respMsg:          "some test response",
			statusCode:       http.StatusPaymentRequired,
			expectedRespBody: `{"error":"some test response"}`,
		},
		{
			name:             "test ok. empty respMsg",
			respMsg:          "",
			statusCode:       http.StatusPaymentRequired,
			expectedRespBody: `{"error":"unknown error"}`,
		},
	}
	//Act
	for _, tc := range testTable {
		w := httptest.NewRecorder()
		respondWithError(w, tc.respMsg, tc.statusCode)
		//Assert
		assert.Equal(t, tc.statusCode, w.Code,
			fmt.Sprintf("%s failed. Wrong status code. Expected %d, got %d",
				tc.name, tc.statusCode, w.Code))
		assert.Equal(t, tc.expectedRespBody, w.Body.String(),
			fmt.Sprintf("%s failed. Wrong response. Expected %s, got %s",
				tc.name, tc.expectedRespBody, w.Body.String()))
	}
}

func TestRespondSuccess(t *testing.T) {
	//Arrange
	testTable := []struct {
		name               string
		url                string
		expectedStatusCode int
		expectedRespBody   string
	}{
		{
			name:               "test ok",
			url:                "https://habr.com/ru/company/avito/blog/658907/",
			expectedStatusCode: http.StatusOK,
			expectedRespBody:   `{"url":"https://habr.com/ru/company/avito/blog/658907/"}`,
		},
		{
			name:               "test ok. empty url",
			url:                "",
			expectedStatusCode: http.StatusInternalServerError,
			expectedRespBody:   `{"error":"server-side error"}`,
		},
	}
	//Act
	for _, tc := range testTable {
		w := httptest.NewRecorder()
		respondSuccess(w, tc.url)
		//Assert
		assert.Equal(t, tc.expectedStatusCode, w.Code,
			fmt.Sprintf("%s failed. Wrong status code. Expected %d, got %d",
				tc.name, tc.expectedStatusCode, w.Code))
		assert.Equal(t, tc.expectedRespBody, w.Body.String(),
			fmt.Sprintf("%s failed. Wrong response. Expected %s, got %s",
				tc.name, tc.expectedRespBody, w.Body.String()))
	}
}
