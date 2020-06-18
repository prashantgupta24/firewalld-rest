package route

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var data1 string
var data2 string
var ipAddr1 string
var ipAddr2 string
var ipAddr3 string

func setup() {
	ipAddr1 = "10.20.30.40"
	ipAddr2 = "20.40.60.80"
	ipAddr3 = "10.50.100.150"
	data1 = `{"ip":"` + ipAddr1 + `","domain":"test.com"}`
	data2 = `{"ip":"` + ipAddr2 + `","domain":"test.com"}`
}

func shutdown() {
	os.Remove("firewalld-rest.db")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestShowAllIPs(t *testing.T) {
	req, err := http.NewRequest("GET", "/ip", nil)
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
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAddIP(t *testing.T) {
	req, err := http.NewRequest("POST", "/ip", strings.NewReader(data1))
	if err != nil {
		t.Fatal(err)
	}
	rr := newRequestRecorder(req, IPAdd)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"meta":null,"data":` + data1 + `}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAddIPDup(t *testing.T) {
	req, err := http.NewRequest("POST", "/ip", strings.NewReader(data1))
	if err != nil {
		t.Fatal(err)
	}
	rr := newRequestRecorder(req, IPAdd)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"error":{"status":400,"title":"ip already exists"}}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestShowIP(t *testing.T) {
	req, err := http.NewRequest("GET", "/ip/"+ipAddr1, nil)
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
	expected := `{"meta":null,"data":` + data1 + `}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestShowIPNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/ip/"+ipAddr3, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/ip/{ip}", IPShow)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"error":{"status":404,"title":"record not found"}}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAddIP2(t *testing.T) {
	req, err := http.NewRequest("POST", "/ip", strings.NewReader(data2))
	if err != nil {
		t.Fatal(err)
	}
	rr := newRequestRecorder(req, IPAdd)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"meta":null,"data":` + data2 + `}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestShowAllIPsAfterAdding(t *testing.T) {
	req, err := http.NewRequest("GET", "/ip", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := newRequestRecorder(req, ShowAllIPs)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//make sure both IPs exist
	if strings.Index(rr.Body.String(), ipAddr1) == -1 || strings.Index(rr.Body.String(), ipAddr2) == -1 {
		t.Errorf("handler returned without required body: \ngot  %v want %v",
			rr.Body.String(), ipAddr1)
	}
}

func TestDeleteIP(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/ip/"+ipAddr1, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/ip/{ip}", IPDelete)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"meta":null,"data":` + data1 + `}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteIPNotFound(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/ip/"+ipAddr1, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/ip/{ip}", IPDelete)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Check the response body is what we expect.
	expected := `{"error":{"status":404,"title":"ip does not exist"}}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: \ngot  %v want %v",
			rr.Body.String(), expected)
	}
}

func TestShowAllIPsAfter(t *testing.T) {
	req, err := http.NewRequest("GET", "/ip", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := newRequestRecorder(req, ShowAllIPs)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//make sure ip doesn't exist
	if strings.Index(rr.Body.String(), ipAddr1) != -1 {
		t.Errorf("handler contained deleted entry: \ngot  %vdeleted %v",
			rr.Body.String(), ipAddr1)
	}
}

// Mocks a handler and returns a httptest.ResponseRecorder
func newRequestRecorder(req *http.Request, fnHandler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fnHandler)
	handler.ServeHTTP(rr, req)
	return rr
}
