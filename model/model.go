package model

import (
	"fmt"
	"log"

	"github.com/firewalld-rest/db"
)

// IPStruct that holds json for ip and domain
type IPStruct struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}

type ipHandler struct {
	filename string
}

//GetIPHandler gets handler for IP
func GetIPHandler() *ipHandler {
	ipHandler := &ipHandler{
		filename: "./firewalld-rest.db",
	}
	return ipHandler
}

func init() {

	ipHandler := GetIPHandler()
	ipStore, err := ipHandler.loadIPStore()
	if err != nil {
		log.Fatal(err)
	}
	if len(ipStore) == 0 {

		//in case you want to store some IPs before hand
		// ipStore["1.2.3.4"] = &IPStruct{
		// 	IP:     "1.2.3.4",
		// 	Domain: "first.com",
		// }
		// ipStore["5.6.7.8"] = &IPStruct{
		// 	IP:     "5.6.7.8",
		// 	Domain: "second",
		// }

		db.Register(ipStore)
		if err := ipHandler.saveIPStore(ipStore); err != nil {
			log.Fatal(err)
		}
	}
}

//GetIP from the db
func (ipHandler *ipHandler) GetIP(ipAddr string) (*IPStruct, error) {
	ipStore, err := ipHandler.loadIPStore()
	if err != nil {
		return nil, err
	}
	ip, ok := ipStore[ipAddr]
	if !ok {
		return nil, fmt.Errorf("record not found")
	}
	return ip, nil
}

//GetAllIPs from the db
func (ipHandler *ipHandler) GetAllIPs() ([]*IPStruct, error) {
	ips := []*IPStruct{}
	ipStore, err := ipHandler.loadIPStore()
	if err != nil {
		return nil, err
	}
	for _, ip := range ipStore {
		ips = append(ips, ip)
	}
	return ips, nil
}

//CheckIPExists checks if IP is in db
func (ipHandler *ipHandler) CheckIPExists(ipAddr string) (bool, error) {
	ipStore, err := ipHandler.loadIPStore()
	if err != nil {
		return false, err
	}
	_, ok := ipStore[ipAddr]
	if ok {
		return true, nil
	}
	return false, nil
}

//AddIP to the db
func (ipHandler *ipHandler) AddIP(ip *IPStruct) error {
	ipStore, err := ipHandler.loadIPStore()
	if err != nil {
		return err
	}
	_, ok := ipStore[ip.IP]
	if ok {
		return fmt.Errorf("ip already exists")
	}
	ipStore[ip.IP] = ip
	if err := ipHandler.saveIPStore(ipStore); err != nil {
		return fmt.Errorf("error while saving to file : %v", err)
	}
	return nil
}

//DeleteIP from the db
func (ipHandler *ipHandler) DeleteIP(ipAddr string) (*IPStruct, error) {
	ipStore, err := ipHandler.loadIPStore()
	if err != nil {
		return nil, err
	}
	ip, ok := ipStore[ipAddr]
	if !ok {
		return nil, fmt.Errorf("record not found")
	}
	delete(ipStore, ipAddr)
	if err := ipHandler.saveIPStore(ipStore); err != nil {
		return nil, fmt.Errorf("error while saving to file : %v", err)
	}
	return ip, nil
}

func (ipHandler *ipHandler) loadIPStore() (map[string]*IPStruct, error) {
	var ipStore = make(map[string]*IPStruct)
	if err := db.Load(ipHandler.filename, &ipStore); err != nil {
		return nil, fmt.Errorf("error while loading from file : %v", err)
	}
	//fmt.Println("ipstore: ", ipStore)
	return ipStore, nil
}

func (ipHandler *ipHandler) saveIPStore(ipStore map[string]*IPStruct) error {
	if err := db.Save(ipHandler.filename, ipStore); err != nil {
		return fmt.Errorf("error while saving to file : %v", err)
	}
	return nil
}
