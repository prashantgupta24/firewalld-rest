package ip

// Instance holds json for ip and domain
type Instance struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}

//Handler interface for handling IP related tasks
type Handler interface {
	GetIP(string) (*Instance, error)
	GetAllIPs() ([]*Instance, error)
	CheckIPExists(ipAddr string) (bool, error)
	AddIP(ip *Instance) (string, error)
	DeleteIP(ipAddr string) (*Instance, string, error)
}
