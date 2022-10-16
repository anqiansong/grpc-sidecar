package proxy

type Config struct {
	GrpcServer Server
}

type Server struct {
	Address string
}
