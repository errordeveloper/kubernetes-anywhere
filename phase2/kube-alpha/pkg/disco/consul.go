package disco

type ConsulDisco struct {
}

func (*ConsulDisco) Init(peers []string) {
}
func (*ConsulDisco) Join(peers []string) {
}
func NewConsulDisco() P2PDiscovery {
	return &ConsulDisco{}
}
