// Package services holds all the services that connects repositories into a business flow
package services

// MessageProcessor is an implementation of the service
type MessageProcessor interface {
	ParseMessage(message string) (string, string, string, float64)
}
