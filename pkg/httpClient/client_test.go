package httpClient

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	title               string
	baseUrl             string
	method              string
	path                string
	defaultHeaders      map[string]string
	body                interface{}
	requestHeaders      map[string]string
	expectedRequestBody string
	expectedHeaders     map[string]string
}

var testCases = []*testCase{
	{
		title:               "POST with body",
		method:              http.MethodPost,
		path:                "/path",
		body:                struct{ Foo string }{"bar"},
		expectedRequestBody: "{\"Foo\":\"bar\"}\n",
	},
	{
		title:               "POST without body",
		method:              http.MethodPost,
		path:                "/path",
		expectedRequestBody: "",
	},
	{
		title:               "PATCH without body",
		method:              http.MethodPatch,
		path:                "/path",
		expectedRequestBody: "",
	},
	{
		title:               "GET with query string",
		method:              http.MethodGet,
		path:                "/path?foo=bar",
		expectedRequestBody: "",
	},
	{
		title:               "GET With headers",
		method:              http.MethodPatch,
		path:                "/path",
		defaultHeaders:      map[string]string{"A": "default", "B": "default"},
		requestHeaders:      map[string]string{"A": "override", "C": "request"},
		expectedHeaders:     map[string]string{"A": "override", "B": "default", "C": "request"},
		expectedRequestBody: "",
	},
}

func TestRequest(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			server := createMockServer(t, tc)
			defer server.Close()

			client := NewClient(server.URL, tc.defaultHeaders)
			_, err := client.Request(tc.method, tc.path, tc.body, tc.requestHeaders)
			assert.Nil(t, err)

		})
	}
}

func TestRequestWithBadUrl(t *testing.T) {
	badUrl := "badProtocol:/host"
	client := NewClient(badUrl, nil)
	_, err := client.Request(http.MethodPatch, "/", nil, map[string]string{})
	assert.NotNil(t, err)
}

func createMockServer(t *testing.T, tc *testCase) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, tc.path, req.URL.String())

		body, _ := ioutil.ReadAll(req.Body)

		if len(tc.expectedRequestBody) > 0 {
			buf := bytes.NewBuffer(body)
			gzipReader, err := gzip.NewReader(buf)
			if err != nil {
				panic(err)
			}
			var gunzippedBuf bytes.Buffer
			_, err = gunzippedBuf.ReadFrom(gzipReader)
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tc.expectedRequestBody, string(gunzippedBuf.Bytes()))
		} else {
			assert.Equal(t, tc.expectedRequestBody, string(body))
		}

		assert.Equal(t, tc.method, req.Method)

		if tc.expectedHeaders != nil {
			for expectedKey, expectedValue := range tc.expectedHeaders {
				actualValue := req.Header.Get(expectedKey)
				assert.Equal(t, expectedValue, actualValue)

			}
		}
		rw.Write([]byte("Halloooooo"))
	}))
}
