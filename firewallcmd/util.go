package firewallcmd

import (
	"fmt"
	"os/exec"
)

//EnableRichRuleForIP enables rich rule for IP access + reloads
//example:
//firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'
func EnableRichRuleForIP(ipAddr string) (string, error) {
	cmd1 := exec.Command(`firewall-cmd`, `--permanent`, "--zone=public", `--add-rich-rule=rule family="ipv4" source address="`+ipAddr+`/32" port protocol="tcp" port="22" accept`)
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
