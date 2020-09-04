package ip

import (
	"os"
	"sync"

	"github.com/prashantgupta24/firewalld-rest/firewallcmd"
)

var once sync.Once

//singleton reference
var handlerInstance *handlerStruct

//handlerStruct for managing IP related tasks
type handlerStruct struct {
	env string
}

//GetHandler gets singleton handler for IP management
func GetHandler() Handler {
	once.Do(func() {
		handlerInstance = &handlerStruct{
			env: os.Getenv("env"),
		}
	})
	return handlerInstance
}

//GetIP from the db
func (handler *handlerStruct) GetIP(ipAddr string) (*Instance, error) {
	ipInstance := &Instance{}
	if handler.env != "local" {
		exists, err := firewallcmd.CheckIPExistsInFirewallRule(ipAddr)
		if err != nil {
			return nil, err
		}
		if exists {
			ipInstance.IP = ipAddr
		}
	}
	return ipInstance, nil
}

//GetAllIPs from the db
func (handler *handlerStruct) GetAllIPs() ([]*Instance, error) {
	ips := []*Instance{}
	if handler.env != "local" {
		ipString, err := firewallcmd.GetIPSInFirewallRule()
		if err != nil {
			return nil, err
		}
		for _, ipAddr := range ipString {
			ipInstance := &Instance{
				IP: ipAddr,
			}
			ips = append(ips, ipInstance)
		}
	}
	return ips, nil
}

//CheckIPExists checks if IP is in db
func (handler *handlerStruct) CheckIPExists(ipAddr string) (bool, error) {
	if handler.env != "local" {
		exists, err := firewallcmd.CheckIPExistsInFirewallRule(ipAddr)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}
	return false, nil
}

//AddIP to the db
func (handler *handlerStruct) AddIP(ip *Instance) (string, error) {
	if handler.env != "local" {
		command, err := firewallcmd.EnableRichRuleForIP(ip.IP)
		if err != nil {
			return "", err
		}
		return command, nil
	}

	return "not supported", nil
}

//DeleteIP from the db
func (handler *handlerStruct) DeleteIP(ipAddr string) (*Instance, string, error) {
	if handler.env != "local" {
		command, err := firewallcmd.DisableRichRuleForIP(ipAddr)
		if err != nil {
			return nil, "", err
		}
		ipInstance := &Instance{
			IP: ipAddr,
		}
		return ipInstance, command, nil
	}
	return nil, "not supported", nil
}
