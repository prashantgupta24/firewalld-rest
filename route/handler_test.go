package route

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var data string
var ipAddr string

func setup() {
	ipAddr = "10.20.30.40"
	data = `{"ip":"` + ipAddr + `","domain":"test.com"}`
}

func shutdown() {
	os.Remove("firewalld-rest-db.tmp")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestShowAllIPs(t *testing.T) {

	defer os.Remove("firewalld-rest-db.tmp")

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

func TestAddIP(t *testing.T) {

	//defer os.Remove("firewalld-rest-db.tmp")

	//data := `{"ip":"10.20.30.40","domain":"test.com"}`

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
}

func TestShowIP(t *testing.T) {

	//data := `{"ip":"10.20.30.40","domain":"test.com"}`

	req, err := http.NewRequest("GET", "/ip/"+ipAddr, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/ip/{ip}", IPShow)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"meta":null,"data":` + data + `}`

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot %v want %v",
			rr.Body.String(), expected)
	}

}

// Mocks a handler and returns a httptest.ResponseRecorder
func newRequestRecorder(req *http.Request, fnHandler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fnHandler)
	handler.ServeHTTP(rr, req)
	return rr
}
