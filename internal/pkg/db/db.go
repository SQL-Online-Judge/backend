package db

type DB interface {
	Connect(connStr string) error
	Close() error
	Ping() error
}
