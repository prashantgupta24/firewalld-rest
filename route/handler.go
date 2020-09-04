package route

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prashantgupta24/firewalld-rest/ip"
)

//Index page
// GET /
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// IPAdd for the Create action
// POST /ip
func IPAdd(w http.ResponseWriter, r *http.Request) {
	ipInstance := &ip.Instance{}
	if err := populateModelFromHandler(w, r, ipInstance); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	ipExists, err := ip.GetHandler().CheckIPExists(ipInstance.IP)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ipExists {
		writeErrorResponse(w, http.StatusBadRequest, "ip already exists")
		return
	}

	command, err := ip.GetHandler().AddIP(ipInstance)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot exec command %v, err : %v", command, err.Error()))
		return
	}

	writeOKResponse(w, ipInstance)
}

// IPShow for the ip Show action
// GET /ip/{ip}
func IPShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ipAddr := vars["ip"]
	ip, err := ip.GetHandler().GetIP(ipAddr)
	if err != nil {
		// No IP found
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	writeOKResponse(w, ip)
}

// ShowAllIPs shows all IPs
// GET /ip
func ShowAllIPs(w http.ResponseWriter, r *http.Request) {
	ips, err := ip.GetHandler().GetAllIPs()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOKResponse(w, ips)
}

// IPDelete for the ip Delete action
// DELETE /ip/{ip}
func IPDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ipAddr := vars["ip"]
	log.Printf("IP to delete %s\n", ipAddr)

	ipExists, err := ip.GetHandler().CheckIPExists(ipAddr)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ipExists {
		writeErrorResponse(w, http.StatusNotFound, "ip does not exist")
		return
	}

	ip, command, err := ip.GetHandler().DeleteIP(ipAddr)
	if err != nil {
		// IP could not be deleted
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot exec command %v, err : %v", command, err.Error()))
		return
	}
	writeOKResponse(w, ip)
}

// Writes the response as a standard JSON response with StatusOK
func writeOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&JSONResponse{Data: m}); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.
		NewEncoder(w).
		Encode(&JSONErrorResponse{Error: &APIError{Status: errorCode, Title: errorMsg}})
}

//Populates a ip from the params in the Handler
func populateModelFromHandler(w http.ResponseWriter, r *http.Request, ip interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, ip); err != nil {
		return err
	}
	return nil
}
