package output

// IptablesStatus 结构体表示iptables的状态
type IptablesStatus struct {
	Enabled     bool // iptables 是否开启
	PingBlocked bool // 是否禁ping
}

// IptablesRule 结构体表示单个iptables规则
type IptablesRule struct {
	Chain  string // 规则所属的链
	Target string // 目标（ACCEPT, DROP等）
	Proto  string // 协议
	Source string // 源IP
	Dest   string // 目标IP
	Port   string // 端口

}
