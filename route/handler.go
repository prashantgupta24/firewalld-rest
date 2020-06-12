package route

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/firewalld-rest/model"
	"github.com/gorilla/mux"
)

//Index page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// IPAdd for the Create action
// POST /ip
func IPAdd(w http.ResponseWriter, r *http.Request) {
	ip := &model.IP{}
	if err := populateModelFromHandler(w, r, ip); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	err := model.AddIP(ip)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	//firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'
	//firewall-cmd -reload

	//command 1
	cmd1 := exec.Command("firewall-cmd", "--zone=public", "--list-all")
	out1, err := cmd1.CombinedOutput()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot exec command %v, err : %v", cmd1.String(), err.Error()))
		return
	}
	fmt.Printf("combined out:\n%s\n", string(out1))

	//command 2
	cmd2 := exec.Command("firewall-cmd", "--zone=private", "--list-all")
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot exec command %v, err : %v", cmd2.String(), err.Error()))
		return
	}
	fmt.Printf("combined out:\n%s\n", string(out2))

	writeOKResponse(w, ip)
}

// IPShow for the ip Show action
// GET /ip/:ip
func IPShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ipAddr := vars["ip"]
	ip, err := model.GetIP(ipAddr)
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
	ips, err := model.GetAllIPs()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOKResponse(w, ips)
}

// IPDelete for the ip Delete action
// DELETE /ip/:ip
func IPDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ipAddr := vars["ip"]
	log.Printf("IP to delete %s\n", ipAddr)
	ip, err := model.DeleteIP(ipAddr)
	if err != nil {
		// No IP found
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	log.Println("IP List before delete")
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

//Populates a model from the params in the Handler
func populateModelFromHandler(w http.ResponseWriter, r *http.Request, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	log.Printf("Response body : %s\n", body)
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}
