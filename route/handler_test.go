package route

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAddIP(t *testing.T) {
	data := `{"ip":"73.223.28.39","domain":"pg.com"}`

	req, err := http.NewRequest("POST", "/ip", strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, IPAdd)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"meta":null,"data":` + data + `}`

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot %vwant %v",
			rr.Body.String(), expected)
	}

	//Check if IP added

	req, err = http.NewRequest("GET", "/ip/73.223.28.39", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = newRequestRecorder(req, IPShow)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = `{"meta":null,"data":` + data + `}`

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot %vwant %v",
			rr.Body.String(), expected)
	}

}

func TestShowAllIPs(t *testing.T) {
	os.Setenv("env", "test")
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/hh", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := newRequestRecorder(req, ShowAllIPs)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"meta":null,"data":[]}`

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot %vwant %v",
			rr.Body.String(), expected)
	}
}

// Mocks a handler and returns a httptest.ResponseRecorder
func newRequestRecorder(req *http.Request, fnHandler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fnHandler)

	handler.ServeHTTP(rr, req)

	// router := httprouter.New()
	// router.Handle(method, strPath, fnHandler)
	// // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	// rr := httptest.NewRecorder()
	// // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// // directly and pass in our Request and ResponseRecorder.
	// router.ServeHTTP(rr, req)
	return rr
}
