package ip

import (
	"fmt"
	"log"
	"sync"

	"github.com/prashantgupta24/firewalld-rest/db"
)

var once sync.Once

//singleton reference
var handlerInstance *handlerStruct

//handlerStruct for managing IP related tasks
type handlerStruct struct {
	db db.Instance
}

//GetHandler gets singleton handler for IP management
func GetHandler() Handler {
	once.Do(func() {
		dbInstance := db.GetFileTypeInstance()
		handlerInstance = &handlerStruct{
			db: dbInstance,
		}
		ipStore, err := handlerInstance.loadIPStore()
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
			handlerInstance.db.Register(ipStore)
			if err := handlerInstance.saveIPStore(ipStore); err != nil {
				log.Fatal(err)
			}
		}
	})
	return handlerInstance
}

//GetIP from the db
func (handler *handlerStruct) GetIP(ipAddr string) (*Instance, error) {
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
func (handler *handlerStruct) GetAllIPs() ([]*Instance, error) {
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
func (handler *handlerStruct) CheckIPExists(ipAddr string) (bool, error) {
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
func (handler *handlerStruct) AddIP(ip *Instance) error {
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
func (handler *handlerStruct) DeleteIP(ipAddr string) (*Instance, error) {
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

func (handler *handlerStruct) loadIPStore() (map[string]*Instance, error) {
	var ipStore = make(map[string]*Instance)
	if err := handler.db.Load(&ipStore); err != nil {
		return nil, fmt.Errorf("error while loading from file : %v", err)
	}
	return ipStore, nil
}

func (handler *handlerStruct) saveIPStore(ipStore map[string]*Instance) error {
	if err := handler.db.Save(ipStore); err != nil {
		return fmt.Errorf("error while saving to file : %v", err)
	}
	return nil
}
