package database

type DbClient interface {
	Health() map[string]string

	CreateProperty() map[string]string
	DeleteProperties() map[string]int64
	CreateProperties() map[string]string
	CountProperties() map[string]int64

	CreateAlert() map[string]string
	// TODO add here all methods to interact with databases. Will be implemented by all DB clientes (so far, only MongoDB)
}
