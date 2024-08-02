package database

import "github.com/venturarome/DaftWatch/internal/model"

type DbClient interface {
	Health() map[string]string

	CreateProperty() map[string]string
	CreateProperties() map[string]string
	DeleteProperties() map[string]int64
	CountProperties() map[string]int64
	//FindPropertiesByListingIds() []model.Property

	CreateAlertForUser(alert model.Alert, user model.User) bool
	DeleteAlerts() map[string]int64

	DeleteUsers() map[string]int64
	// TODO add here all methods to interact with databases. Will be implemented by all DB clientes (so far, only MongoDB)
}
