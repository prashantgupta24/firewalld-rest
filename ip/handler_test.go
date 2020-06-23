package ip

import (
	"log"
	"os"
	"os/exec"
	"strings"
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

func TestGetAllIPsFileError(t *testing.T) {
	changeFilePermission("100")
	_, err := handler.GetAllIPs()
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && strings.Index(err.Error(), "permission denied") == -1 {
		t.Errorf("should have received permission error, instead got : %v", err)
	}
	changeFilePermission("644")
}

func TestGetAllIPs(t *testing.T) {
	ips, err := handler.GetAllIPs()
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
	if len(ips) != 0 {
		t.Errorf("should have been an empty list, instead got : %v", len(ips))
	}
}

func TestAddIPFileError(t *testing.T) {
	changeFilePermission("100")
	err := handler.AddIP(ipInstance)
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && strings.Index(err.Error(), "permission denied") == -1 {
		t.Errorf("should have received permission error, instead got : %v", err)
	}

	changeFilePermission("500")
	err = handler.AddIP(ipInstance)
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && strings.Index(err.Error(), "permission denied") == -1 {
		t.Errorf("should have received permission error, instead got : %v", err)
	}
	changeFilePermission("644")
}
func TestAddIP(t *testing.T) {
	err := handler.AddIP(ipInstance)
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
}

func TestAddIPDup(t *testing.T) {
	err := handler.AddIP(ipInstance)
	if err == nil {
		t.Errorf("should have errored for duplicate IP")
	}
}

func TestGetAllIPsAfterAdd(t *testing.T) {
	ips, err := handler.GetAllIPs()
	if err != nil {
		t.Errorf("should not have errored, err : %v", err)
	}
	if len(ips) == 0 {
		t.Errorf("should have included %v , instead got : %v", ipAddr, ips)
	}
}

func TestCheckIPExistsFileError(t *testing.T) {
	changeFilePermission("100")
	_, err := handler.CheckIPExists(ipAddr)
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && strings.Index(err.Error(), "permission denied") == -1 {
		t.Errorf("should have received permission error, instead got : %v", err)
	}
	changeFilePermission("644")
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

func TestGetIPFileError(t *testing.T) {
	changeFilePermission("100")
	_, err := handler.GetIP(ipAddr)
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && strings.Index(err.Error(), "permission denied") == -1 {
		t.Errorf("should have received permission error, instead got : %v", err)
	}
	changeFilePermission("644")
}

func TestGetInvalidIP(t *testing.T) {
	_, err := handler.GetIP("invalid_ip")
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && err.Error() != "record not found" {
		t.Errorf("record should not have been found, instead got : %v", err)
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

func TestDeleteIPFileError(t *testing.T) {
	changeFilePermission("100")
	_, err := handler.DeleteIP(ipAddr)
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && strings.Index(err.Error(), "permission denied") == -1 {
		t.Errorf("should have received permission error, instead got : %v", err)
	}

	changeFilePermission("500")
	_, err = handler.DeleteIP(ipAddr)
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && strings.Index(err.Error(), "permission denied") == -1 {
		t.Errorf("should have received permission error, instead got : %v", err)
	}
	changeFilePermission("644")
}

func TestDeleteInvalidIP(t *testing.T) {
	_, err := handler.DeleteIP("invalid_ip")
	if err == nil {
		t.Errorf("should have errored")
	}
	if err != nil && err.Error() != "record not found" {
		t.Errorf("record should not have been found, instead got : %v", err)
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

func changeFilePermission(permission string) {
	cmd := exec.Command("chmod", permission, "firewalld-rest.db")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not change permission of file, err : %v", err)
	}
	// cmd1 := exec.Command("ls", "-la")
	// o1, _ := cmd1.CombinedOutput()
	// fmt.Println("cmd1 : ", string(o1))
}
