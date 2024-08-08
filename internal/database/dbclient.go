package database

import "github.com/venturarome/DaftWatch/internal/model"

type DbClient interface {
	Health() map[string]string

	CreateUser(user model.User) map[string]interface{}
	ListAlertsForUser(user model.User) []model.Alert
	DeleteUsers() map[string]int64

	AddSubscriberToAlert(alert model.Alert, user model.User) map[string]interface{}
	RemoveSubscriberFromAlert(alert model.Alert, user model.User) bool
	SetPropertiesToAlert(alert model.Alert, properties []model.Property) map[string]interface{}
	DeleteAlerts() map[string]int64

	CreateProperty(property model.Property) map[string]interface{}
	CreateProperties(properties []model.Property) map[string]interface{}
	CountProperties() map[string]int64
	DeleteProperties() map[string]int64

	// TODO add here all methods to interact with databases. Will be implemented by all DB clientes (so far, only MongoDB)
}
