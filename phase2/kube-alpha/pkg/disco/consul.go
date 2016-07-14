package disco

type ConsulDisco struct {
}

func (*ConsulDisco) Init() {
}
func (*ConsulDisco) Bootstrap(peers []IPv4Address) {
}
func NewConsulDisco() P2PDiscovery {
	return &ConsulDisco{}
}
