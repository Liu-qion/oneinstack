package output

type FirewallStatus struct {
	Enabled bool `json:"enabled"`
}

type FirewallPort struct {
	Protocol  string `json:"protocol"`
	Port      string `json:"port"`
	State     string `json:"state"`
	Policy    string `json:"policy"`
	Direction string `json:"direction"`
	Source    string `json:"source"`
}

type FirewallPorts struct {
	Ports []FirewallPort `json:"ports"`
}
