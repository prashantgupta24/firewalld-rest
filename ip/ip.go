package ip

import (
	"fmt"
	"log"
	"os"

	"github.com/firewalld-rest/db"
)

// Instance holds json for ip and domain
type Instance struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}

//handler for managing IP related tasks
type handler struct {
	db db.Instance
}

//GetHandler gets handler for IP
func GetHandler() handler {

	path := os.Getenv("FIREWALLD_REST_DB_PATH")
	fmt.Println("FIREWALLD_REST_DB_PATH : ", path)
	if path != "" {
		path = parsePath(path)
		path += "/firewalld-rest.db"
	} else {
		path = "./firewalld-rest.db"
	}
	handler := handler{
		db: &db.FileType{
			Path: path,
		},
	}
	return handler
}

func init() {

	handler := GetHandler()
	ipStore, err := handler.loadIPStore()
	if err != nil {
		log.Fatal(err)
	}
	if len(ipStore) == 0 {

		//in case you want to store some IPs before hand
		// ipStore["1.2.3.4"] = &Instance{
		// 	IP:     "1.2.3.4",
		// 	Domain: "first.com",
		// }
		// ipStore["5.6.7.8"] = &Instance{
		// 	IP:     "5.6.7.8",
		// 	Domain: "second",
		// }

		handler.db.Register(ipStore)
		if err := handler.saveIPStore(ipStore); err != nil {
			log.Fatal(err)
		}
	}
}

//GetIP from the db
func (handler handler) GetIP(ipAddr string) (*Instance, error) {
	ipStore, err := handler.loadIPStore()
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
func (handler handler) GetAllIPs() ([]*Instance, error) {
	ips := []*Instance{}
	ipStore, err := handler.loadIPStore()
	if err != nil {
		return nil, err
	}
	for _, ip := range ipStore {
		ips = append(ips, ip)
	}
	return ips, nil
}

//CheckIPExists checks if IP is in db
func (handler handler) CheckIPExists(ipAddr string) (bool, error) {
	ipStore, err := handler.loadIPStore()
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
func (handler handler) AddIP(ip *Instance) error {
	ipStore, err := handler.loadIPStore()
	if err != nil {
		return err
	}
	_, ok := ipStore[ip.IP]
	if ok {
		return fmt.Errorf("ip already exists")
	}
	ipStore[ip.IP] = ip
	if err := handler.saveIPStore(ipStore); err != nil {
		return fmt.Errorf("error while saving to file : %v", err)
	}
	return nil
}

//DeleteIP from the db
func (handler handler) DeleteIP(ipAddr string) (*Instance, error) {
	ipStore, err := handler.loadIPStore()
	if err != nil {
		return nil, err
	}
	ip, ok := ipStore[ipAddr]
	if !ok {
		return nil, fmt.Errorf("record not found")
	}
	delete(ipStore, ipAddr)
	if err := handler.saveIPStore(ipStore); err != nil {
		return nil, fmt.Errorf("error while saving to file : %v", err)
	}
	return ip, nil
}

func (handler handler) loadIPStore() (map[string]*Instance, error) {
	var ipStore = make(map[string]*Instance)
	if err := handler.db.Load(&ipStore); err != nil {
		return nil, fmt.Errorf("error while loading from file : %v", err)
	}
	//fmt.Println("ipstore: ", ipStore)
	return ipStore, nil
}

func (handler handler) saveIPStore(ipStore map[string]*Instance) error {
	if err := handler.db.Save(ipStore); err != nil {
		return fmt.Errorf("error while saving to file : %v", err)
	}
	return nil
}

func parsePath(path string) string {
	lastChar := path[len(path)-1:]

	if lastChar == "/" {
		path = path[:len(path)-1]
	}
	return path
}
