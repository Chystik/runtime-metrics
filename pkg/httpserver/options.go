package httpserver

type Options func(*Server)

func Address(addr string) Options {
	return func(s *Server) {
		s.Addr = addr
	}
}
