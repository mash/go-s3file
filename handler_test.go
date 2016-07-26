package s3file

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`ok`))
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	req.Host = "example.com"
	return req
}

func testHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	mux.HandleFunc("/.well-known/acme-challenge/", Handler)
	return mux
}

func TestHandler(t *testing.T) {
	tests := []struct {
		Method string
		Path   string
		Code   int
	}{
		{
			Method: "GET",
			Path:   "/",
			Code:   200,
		},
		{
			Method: "GET",
			Path:   "/.well-known/acme-challenge/404",
			Code:   404,
		},
		{
			Method: "GET",
			Path:   "/.well-known/acme-challenge/200",
			Code:   200,
		},
	}
	for _, test := range tests {
		t.Logf("Testing %s %s", test.Method, test.Path)
		writer := httptest.NewRecorder()
		testHandler().ServeHTTP(writer, newRequest(test.Method, test.Path))
		if writer.Code != test.Code {
			t.Errorf("expected %d but got %d", test.Code, writer.Code)
		}
	}
}