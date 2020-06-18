package firewallcmd

import "os/exec"

//EnableRichRuleForIP enables rich rule for IP access
//example:
//firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'
func EnableRichRuleForIP(ipAddr string) (string, string, error) {
	cmd := exec.Command(`firewall-cmd`, `--permanent`, "--zone=public", `--add-rich-rule=rule family="ipv4" source address="`+ipAddr+`/32" port protocol="tcp" port="22" accept`)
	//uncomment for debugging
	// for _, v := range cmd1.Args {
	// 	fmt.Println(v)
	// }
	output, err := cmd.CombinedOutput()
	return cmd.String(), string(output), err
}

//DisableRichRuleForIP disables rich rule for IP access
func DisableRichRuleForIP(ipAddr string) (string, string, error) {
	cmd := exec.Command(`firewall-cmd`, `--permanent`, "--zone=public", `--remove-rich-rule=rule family="ipv4" source address="`+ipAddr+`/32" port protocol="tcp" port="22" accept`)
	output, err := cmd.CombinedOutput()
	return cmd.String(), string(output), err
}

//Reload reloads firewall for setting to take effect
func Reload() (string, string, error) {
	cmd := exec.Command("firewall-cmd", "--reload")
	output, err := cmd.CombinedOutput()
	return cmd.String(), string(output), err
}
