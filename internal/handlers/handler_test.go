package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUrlFromRequest(t *testing.T) {
	//Arrange
	testCases := []struct {
		name             string
		httpMethod       string
		reqBody          io.Reader
		WantUrl          string
		WantErr          error
		WantHTTPStatus   int
		WantJsonResponse string
	}{
		{
			name:             "POST with empty body",
			httpMethod:       http.MethodPost,
			reqBody:          strings.NewReader(""),
			WantUrl:          "",
			WantErr:          fmt.Errorf("empty body"),
			WantHTTPStatus:   http.StatusBadRequest,
			WantJsonResponse: `{"error":"empty body"}`,
		},
		{
			name:             "POST with empty JSON",
			httpMethod:       http.MethodPost,
			reqBody:          strings.NewReader("{}"),
			WantUrl:          "",
			WantErr:          fmt.Errorf("incorrect data in body"),
			WantHTTPStatus:   http.StatusBadRequest,
			WantJsonResponse: `{"error":"incorrect data in body"}`,
		},
		{
			name:             "POST with wrong data in body 1",
			httpMethod:       http.MethodPost,
			reqBody:          strings.NewReader(`{"url": aboba.com}`),
			WantUrl:          "",
			WantErr:          fmt.Errorf("unable recognize data in body: invalid character 'a' looking for beginning of value"),
			WantHTTPStatus:   http.StatusBadRequest,
			WantJsonResponse: `{"error":"unable recognize data in body"}`,
		},
		{
			name:             "POST with correct request",
			httpMethod:       http.MethodPost,
			reqBody:          strings.NewReader(`{"url": "https://www.ozon.ru/product/falloimitator-hello-kiki-dildo-m-fioletovyy-lw-d-sl-18-5-sm-faloimitator-dildo-na-prisoske-fallos-478864885/"}`),
			WantUrl:          "https://www.ozon.ru/product/falloimitator-hello-kiki-dildo-m-fioletovyy-lw-d-sl-18-5-sm-faloimitator-dildo-na-prisoske-fallos-478864885/",
			WantErr:          nil,
			WantHTTPStatus:   http.StatusOK,
			WantJsonResponse: "",
		},
	}

	//Act
	for _, tc := range testCases {
		rec := httptest.NewRecorder()
		request := httptest.NewRequest(tc.httpMethod, "/", tc.reqBody)
		url, err := getUrlFromRequest(rec, request)
		//Assert
		if url != tc.WantUrl {
			t.Errorf("getUrlFromRequest returns incorrect url: want %s, got %s",
				tc.WantUrl, url)
		}
		//TODO define errors in package and test via errors.IS() / errors.AS()
		if err != nil {
			if err.Error() != tc.WantErr.Error() {
				t.Errorf("getUrlFromRequest returns incorrect error: want \"%s\", got \"%s\"",
					tc.WantErr, err)
			}
		}
		if rec.Code != tc.WantHTTPStatus {
			t.Errorf("getUrlFromRequest returns incorrect status code: want %d, got %d",
				tc.WantHTTPStatus, rec.Code)
		}
		if rec.Body.String() != tc.WantJsonResponse {
			t.Errorf("getUrlFromRequest returns incorrect response: want %s, got %s",
				tc.WantJsonResponse, rec.Body.String())
		}
	}
}

func TestProcessShort(t *testing.T) {

}

func TestProcessLong(t *testing.T) {

}