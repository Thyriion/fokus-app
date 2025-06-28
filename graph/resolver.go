package graph

import (
	"fokus-app/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	focusareaService *db.FocusareaService
	sessionService   *db.SessionService
}

func NewResolver(database *db.Database) *Resolver {
	return &Resolver{
		focusareaService: db.NewFocusareaService(database),
		sessionService:   db.NewSessionService(database),
	}
}
