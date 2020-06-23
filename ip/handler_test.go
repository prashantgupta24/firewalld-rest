package ip

import (
	"os"
	"testing"
)

var handler Handler
var ipAddr string
var ipInstance *Instance

func setup() {
	handler = GetHandler()
	ipAddr = "10.20.30.40"
	ipInstance = &Instance{
		IP:     ipAddr,
		Domain: "test",
	}
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

func TestGetHandler(t *testing.T) {
	if handler == nil {
		t.Errorf("handler should not be nil")
	}
}

func TestGetAllIPs(t *testing.T) {
	ips, err := handler.GetAllIPs()
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
	if len(ips) != 0 {
		t.Errorf("should have been an empty list, instead got : %v", ips)
	}
}

func TestAddIP(t *testing.T) {
	err := handler.AddIP(ipInstance)
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
}

func TestCheckIPExists(t *testing.T) {
	ipExists, err := handler.CheckIPExists(ipAddr)
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
	if !ipExists {
		t.Errorf("ip %v should exist", ipAddr)
	}
}

func TestGetIP(t *testing.T) {
	ipRecd, err := handler.GetIP(ipAddr)
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
	if ipRecd.IP != ipInstance.IP {
		t.Errorf("ip should be same, got %v want %v", ipRecd.IP, ipInstance.IP)
	}
}

func TestDeleteIP(t *testing.T) {
	ipDeleted, err := handler.DeleteIP(ipAddr)
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
	if ipAddr != ipDeleted.IP {
		t.Errorf("ip %v should be the same", ipAddr)
	}
	ipExists, err := handler.CheckIPExists(ipAddr)
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
	if ipExists {
		t.Errorf("ip %v should be deleted", ipAddr)
	}
}
