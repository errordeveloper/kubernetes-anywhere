package disco

type ConsulDisco struct {
}

func (*ConsulDisco) Bootstrap(peers []string) {
}
func NewConsulDisco() P2PDiscovery {
	return &ConsulDisco{}
}
