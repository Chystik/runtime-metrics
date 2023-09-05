package storage

type MetricsStorage interface {
	Read() error
	Write() error
	CloseFile() error
}
