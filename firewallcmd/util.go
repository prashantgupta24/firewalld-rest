package firewallcmd

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

//EnableRichRuleForIP enables rich rule for IP access + reloads
//example:
//firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'
func EnableRichRuleForIP(ipAddr string) (string, error) {

	//check if valid ipv4 address
	if !isValidIpv4(ipAddr) {
		return "", fmt.Errorf("not a valid IPv4 address : %v", ipAddr)
	}
	cmd1 := exec.Command(`firewall-cmd`, `--permanent`, "--zone=public", `--add-rich-rule=`+createRichRule(ipAddr))
	//uncomment for debugging
	// for _, v := range cmd1.Args {
	// 	fmt.Println(v)
	// }
	output1, err1 := cmd1.CombinedOutput()
	if err1 != nil {
		return cmd1.String(), err1
	}
	fmt.Printf("rich rule added successfully for ip %v : %v", ipAddr, string(output1))

	cmd2, output2, err2 := reload()
	if err2 != nil {
		return cmd2.String(), err2
	}
	fmt.Printf("firewalld reloaded successfully : %v", string(output2))
	return "", nil
}

//DisableRichRuleForIP disables rich rule for IP access + reloads
func DisableRichRuleForIP(ipAddr string) (string, error) {
	cmd1 := exec.Command(`firewall-cmd`, `--permanent`, "--zone=public", `--remove-rich-rule=rule family="ipv4" source address="`+ipAddr+`/32" port protocol="tcp" port="22" accept`)
	output1, err1 := cmd1.CombinedOutput()
	if err1 != nil {
		return cmd1.String(), err1
	}
	fmt.Printf("rich rule deleted successfully for ip %v : %v", ipAddr, string(output1))

	cmd2, output2, err2 := reload()
	if err2 != nil {
		return cmd2.String(), err2
	}
	fmt.Printf("firewalld reloaded successfully : %v", string(output2))
	return "", nil
}

//reload reloads firewall for setting to take effect
func reload() (*exec.Cmd, []byte, error) {
	cmd := exec.Command("firewall-cmd", "--reload")
	output, err := cmd.CombinedOutput()
	return cmd, output, err
}

//GetIPSInFirewall gets IPs currently in firewall
func GetIPSInFirewallRule() ([]string, error) {

	var ipsInFirewall []string
	cmd := exec.Command("firewall-cmd", "--zone=public", "--list-rich-rules")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error while fetching IPs from firewall-cmd, err: %v", err)
	}

	stringToSearch := "address=\""
	richRuleLines := strings.Split(string(output), "\n")

	for _, rule := range richRuleLines {
		r, _ := regexp.Compile(stringToSearch + "[0-9.]*")
		ipAddr := strings.TrimPrefix(r.FindString(rule), stringToSearch)
		if isValidIpv4(ipAddr) {
			ipsInFirewall = append(ipsInFirewall, ipAddr)
			//fmt.Println(ipAddr)
		}
	}
	return ipsInFirewall, err
}

//CheckIPExistsInFirewallRule checks if rich rule exists with IP
func CheckIPExistsInFirewallRule(ipAddr string) (bool, error) {
	cmd := exec.Command(`firewall-cmd`, "--zone=public", `--query-rich-rule=`+createRichRule(ipAddr))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	strOutput := string(output)
	if strOutput == "yes" {
		return true, nil
	} else if strOutput == "no" {
		return false, nil
	}
	return false, nil
}

func createRichRule(ipAddr string) string {
	richRule := `rule family="ipv4" source address="` + ipAddr + `/32" port protocol="tcp" port="22" accept`
	return richRule
}
func isValidIpv4(host string) bool {
	return net.ParseIP(host) != nil
}
