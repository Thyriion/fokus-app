package db

import (
	"database/sql"
	"fmt"
	"fokus-app/graph/model"
	"strconv"
	"time"
)

type SessionService struct {
	db *Database
}

func NewSessionService(db *Database) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) CreateSession(focusareaID string) (*model.Session, error) {
	// Verify the focusarea exists
	focusareaService := NewFocusareaService(s.db)
	_, err := focusareaService.GetFocusareaByID(focusareaID)
	if err != nil {
		return nil, fmt.Errorf("invalid focusarea: %w", err)
	}

	startTime := time.Now().Format(time.RFC3339)
	query := `INSERT INTO sessions (focusarea_id, start_time) VALUES (?, ?)`

	result, err := s.db.conn.Exec(query, focusareaID, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &model.Session{
		ID:              strconv.FormatInt(id, 10),
		FocusareaID:     focusareaID,
		Start:           startTime,
		End:             nil,
		DurationMinutes: nil,
	}, nil
}

func (s *SessionService) GetAllSessions() ([]*model.Session, error) {
	query := `SELECT id, focusarea_id, start_time, end_time, duration_minutes FROM sessions ORDER BY created_at DESC`

	rows, err := s.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*model.Session
	for rows.Next() {
		var id, focusareaID int64
		var startTime string
		var endTime sql.NullString
		var durationMinutes sql.NullInt32

		if err := rows.Scan(&id, &focusareaID, &startTime, &endTime, &durationMinutes); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		session := &model.Session{
			ID:          strconv.FormatInt(id, 10),
			FocusareaID: strconv.FormatInt(focusareaID, 10),
			Start:       startTime,
		}

		if endTime.Valid {
			session.End = &endTime.String
		}

		if durationMinutes.Valid {
			session.DurationMinutes = &durationMinutes.Int32
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *SessionService) GetActiveSession() (*model.Session, error) {
	query := `SELECT id, focusarea_id, start_time FROM sessions WHERE end_time IS NULL ORDER BY created_at DESC LIMIT 1`

	var id, focusareaID int64
	var startTime string

	err := s.db.conn.QueryRow(query).Scan(&id, &focusareaID, &startTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active session found")
		}
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}

	return &model.Session{
		ID:              strconv.FormatInt(id, 10),
		FocusareaID:     strconv.FormatInt(focusareaID, 10),
		Start:           startTime,
		End:             nil,
		DurationMinutes: nil,
	}, nil
}

func (s *SessionService) StopSession() (*model.Session, error) {
	// Get the active session
	activeSession, err := s.GetActiveSession()
	if err != nil {
		return nil, err
	}

	// Calculate duration
	endTime := time.Now()
	startTime, err := time.Parse(time.RFC3339, activeSession.Start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}

	duration := int32(endTime.Sub(startTime).Minutes())
	endTimeStr := endTime.Format(time.RFC3339)

	// Update the session
	query := `UPDATE sessions SET end_time = ?, duration_minutes = ? WHERE id = ?`
	_, err = s.db.conn.Exec(query, endTimeStr, duration, activeSession.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to stop session: %w", err)
	}

	// Return updated session
	activeSession.End = &endTimeStr
	activeSession.DurationMinutes = &duration

	return activeSession, nil
}
