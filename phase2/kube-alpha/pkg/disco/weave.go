package disco

type WeaveDisco struct {
}

func (*WeaveDisco) Bootstrap(peers []IPv4Address) {
	// weave launch $peers
}
func NewWeaveDisco() P2PDiscovery {
	return &WeaveDisco{}
}
