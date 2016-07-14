package disco

type WeaveDisco struct {
}

func (*WeaveDisco) Init() {
}
func (*WeaveDisco) Bootstrap(peers []IPv4Address) {
}
func NewWeaveDisco() P2PDiscovery {
	return &WeaveDisco{}
}
