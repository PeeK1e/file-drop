package environment

type httpServerOptions struct {
	listenAddress string
	port          string
}

func (s *httpServerOptions) String() string {
	return s.listenAddress + ":" + s.port
}
