package graph

import "fokus-app/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	focusareas []*model.Focusarea
	sessions   []*model.Session
}
