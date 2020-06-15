package model

import (
	"fmt"
	"log"

	"github.com/firewalld-rest/db"
)

const (
	filename = "/tmp/firewalld-rest-db.tmp"
)

// IP struct that holds json for ip and domain
type IP struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}

func init() {

	ipStore, err := getIPStoreFromDB()
	if err != nil {
		log.Fatal(err)
	}
	if len(ipStore) == 0 {

		//in case you want to store some IPs before hand
		// ipStore["1.2.3.4"] = &IP{
		// 	IP:     "1.2.3.4",
		// 	Domain: "first.com",
		// }
		// ipStore["5.6.7.8"] = &IP{
		// 	IP:     "5.6.7.8",
		// 	Domain: "second",
		// }

		db.Register(ipStore)
		if err := db.Save(filename, ipStore); err != nil {
			log.Fatal(err)
		}
	}
}

//GetIP from the db
func GetIP(ipAddr string) (*IP, error) {
	ipStore, err := getIPStoreFromDB()
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
func GetAllIPs() ([]*IP, error) {
	ips := []*IP{}
	ipStore, err := getIPStoreFromDB()
	if err != nil {
		return nil, err
	}
	for _, ip := range ipStore {
		ips = append(ips, ip)
	}
	return ips, nil
}

//CheckIPExists checks if IP is in db
func CheckIPExists(ipAddr string) (bool, error) {
	ipStore, err := getIPStoreFromDB()
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
func AddIP(ip *IP) error {
	ipStore, err := getIPStoreFromDB()
	if err != nil {
		return err
	}
	_, ok := ipStore[ip.IP]
	if ok {
		return fmt.Errorf("ip already exists")
	}
	ipStore[ip.IP] = ip
	if err := db.Save(filename, ipStore); err != nil {
		return fmt.Errorf("error while saving to file : %v", err)
	}
	return nil
}

//DeleteIP from the db
func DeleteIP(ipAddr string) (*IP, error) {
	ipStore, err := getIPStoreFromDB()
	if err != nil {
		return nil, err
	}
	ip, ok := ipStore[ipAddr]
	if !ok {
		return nil, fmt.Errorf("record not found")
	}
	delete(ipStore, ipAddr)
	if err := db.Save(filename, ipStore); err != nil {
		return nil, fmt.Errorf("error while saving to file : %v", err)
	}
	return ip, nil

}

func getIPStoreFromDB() (map[string]*IP, error) {
	var ipStore = make(map[string]*IP)
	if err := db.Load(filename, &ipStore); err != nil {
		return nil, fmt.Errorf("error while loading from file : %v", err)
	}
	//fmt.Println("ipstore: ", ipStore)
	return ipStore, nil
}
