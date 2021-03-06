package handlers

import (
	mock_services "URL_shortener_2/internal/services/moks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
			t.Errorf("Test %s. getUrlFromRequest returns incorrect url: want %s, got %s",
				tc.name, tc.WantUrl, url)
		}
		//TODO define errors in package and test via errors.IS() / errors.AS()
		if err != nil {
			if err.Error() != tc.WantErr.Error() {
				t.Errorf("Test %s. getUrlFromRequest returns incorrect error: want \"%s\", got \"%s\"",
					tc.name, tc.WantErr, err)
			}
		}
		if rec.Code != tc.WantHTTPStatus {
			t.Errorf("Test %s. getUrlFromRequest returns incorrect status code: want %d, got %d",
				tc.name, tc.WantHTTPStatus, rec.Code)
		}
		if rec.Body.String() != tc.WantJsonResponse {
			t.Errorf("Test %s. getUrlFromRequest returns incorrect response: want %s, got %s",
				tc.name, tc.WantJsonResponse, rec.Body.String())
		}
	}
}

func TestHandler_ServeHTTP(t *testing.T) {
	//Arrange
	type mockBehaviour func(s *mock_services.MockService, url string)

	testTable := []struct {
		//name of test for debug messages
		name string
		//what we get from client
		requestBody string
		//what getUrlFromRequest should retrieve from request's body
		url string
		//method in request's header (POST for processLong, GET for processShort)
		method string
		//...
		mockBehaviour mockBehaviour
		//What status code do we expect in response's header
		expectedStatusCode int
		//What content do we expect in response's body
		expectedResponseBody string
	}{
		{
			name:        "ok_shortUrl",
			requestBody: `{"url":"RandS10Len"}`,
			url:         "RandS10Len",
			method:      http.MethodGet,
			mockBehaviour: func(s *mock_services.MockService, url string) {
				s.EXPECT().Get(url).Return("https://habr.com/ru/post/425025/", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"url":"https://habr.com/ru/post/425025/"}`,
		},
		{
			name:        "ok_longUrl_that_already_in_storage",
			requestBody: `{"url":"https://youtu.be/Mvw5fbHGJFw"}`,
			method:      http.MethodPost,
			url:         "https://youtu.be/Mvw5fbHGJFw",
			mockBehaviour: func(s *mock_services.MockService, url string) {
				s.EXPECT().Get(url).Return("randSLen10", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"url":"randSLen10"}`,
		},
		//{
		//	name:        "ok_longUrl_that_wasnt_in_storage",
		//	requestBody: `{"url":"https://youtu.be/Mvw5fbHGJFw"}`,
		//	method:      http.MethodPost,
		//	url:         "https://youtu.be/Mvw5fbHGJFw",
		//	mockBehaviour: func(s *mock_services.MockService, url string) {
		//		gomock.InOrder(
		//			s.EXPECT().Get(url),
		//			s.EXPECT().Save(url).Return("randSLen10", nil),
		//		)
		//	},
		//	expectedStatusCode:   http.StatusCreated,
		//	expectedResponseBody: `{"url":"randSLen10"}`,
		//},
		{
			name:                 "not_ok_unsupported method",
			requestBody:          `{"url":"https://youtu.be/Mvw5fbHGJFw"}`,
			method:               http.MethodPatch,
			url:                  "https://youtu.be/Mvw5fbHGJFw",
			mockBehaviour:        func(s *mock_services.MockService, url string) {},
			expectedStatusCode:   http.StatusMethodNotAllowed,
			expectedResponseBody: `{"error":"unsupported method"}`,
		},
	}
	//Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			//gomock's requirement
			c := gomock.NewController(t)
			defer c.Finish()

			//create mocked service
			service := mock_services.NewMockService(c)
			//run function ?FOR WHAT?
			tc.mockBehaviour(service, tc.url)

			//Real handler based on mocked service
			handler := New(service)
			//dummy ResponseWriter and dummy request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.requestBody))
			//Real handler's ServeHTTP call with dummy params
			handler.ServeHTTP(w, req)

			//Assertion
			assert.Equal(t, tc.expectedStatusCode, w.Code,
				fmt.Sprintf("%s test fail. Wrong status code. Expected %d, got %d",
					tc.name, tc.expectedStatusCode, w.Code))
			assert.Equal(t, tc.expectedResponseBody, w.Body.String(),
				fmt.Sprintf("%s test fail. Wrong response. Expected %s, got %s",
					tc.name, tc.expectedResponseBody, w.Body.String()))
		})
	}
}
