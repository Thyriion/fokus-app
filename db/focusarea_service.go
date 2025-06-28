package db

import (
	"database/sql"
	"fmt"
	"fokus-app/graph/model"
	"strconv"
)

type FocusareaService struct {
	db *Database
}

func NewFocusareaService(db *Database) *FocusareaService {
	return &FocusareaService{db: db}
}

func (s *FocusareaService) CreateFocusarea(name string, deadline *string) (*model.Focusarea, error) {
	query := `INSERT INTO focusareas (name, deadline) VALUES (?, ?)`

	result, err := s.db.conn.Exec(query, name, deadline)
	if err != nil {
		return nil, fmt.Errorf("failed to create focusarea: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &model.Focusarea{
		ID:       strconv.FormatInt(id, 10),
		Name:     name,
		Deadline: deadline,
	}, nil
}

func (s *FocusareaService) GetAllFocusareas() ([]*model.Focusarea, error) {
	query := `SELECT id, name, deadline FROM focusareas ORDER BY created_at DESC`

	rows, err := s.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get focusareas: %w", err)
	}
	defer rows.Close()

	var focusareas []*model.Focusarea
	for rows.Next() {
		var id int64
		var name string
		var deadline sql.NullString

		if err := rows.Scan(&id, &name, &deadline); err != nil {
			return nil, fmt.Errorf("failed to scan focusarea: %w", err)
		}

		focusarea := &model.Focusarea{
			ID:   strconv.FormatInt(id, 10),
			Name: name,
		}

		if deadline.Valid {
			focusarea.Deadline = &deadline.String
		}

		focusareas = append(focusareas, focusarea)
	}

	return focusareas, nil
}

func (s *FocusareaService) GetFocusareaByID(id string) (*model.Focusarea, error) {
	query := `SELECT id, name, deadline FROM focusareas WHERE id = ?`

	var focusareaID int64
	var name string
	var deadline sql.NullString

	err := s.db.conn.QueryRow(query, id).Scan(&focusareaID, &name, &deadline)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("focusarea not found")
		}
		return nil, fmt.Errorf("failed to get focusarea: %w", err)
	}

	focusarea := &model.Focusarea{
		ID:   strconv.FormatInt(focusareaID, 10),
		Name: name,
	}

	if deadline.Valid {
		focusarea.Deadline = &deadline.String
	}

	return focusarea, nil
}
