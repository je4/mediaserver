package service

type ShutdownService interface {
	Stop()
	GracefulStop()
}
