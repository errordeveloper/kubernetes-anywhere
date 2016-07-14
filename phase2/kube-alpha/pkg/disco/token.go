package disco

type TokenDisco struct {
}

func (*TokenDisco) Init() string {
	// TODO implementation fill this in please!
	return "<your token here>"
}
func (*TokenDisco) Bootstrap(token string) {
}
func NewTokenDisco() TokenDiscovery {
	return &TokenDisco{}
}
