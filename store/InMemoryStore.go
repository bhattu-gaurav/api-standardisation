package store

import (
	"api-standardisation/tsp-output/notesapi"
	"context"

	"github.com/google/uuid"
)

type InMemoryStore struct {
	users map[string]*notesapi.User
	notes map[string]*notesapi.Note
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: make(map[string]*notesapi.User),
		notes: make(map[string]*notesapi.Note),
	}
}
func (s *InMemoryStore) RegisterUser(ctx context.Context, user *notesapi.User) (*notesapi.User, error) {

	for _, existingUser := range s.users {
		if existingUser.Username == user.Username {
			return nil, notesapi.ErrNoteNotFound
		}
	}

	// Generate ID if not provided
	if user.Id == "" {
		user.Id = uuid.New().String()
	}

	// Store user
	s.users[user.Id] = user
	return user, nil
}

// LoginUser authenticates a user
func (s *InMemoryStore) LoginUser(ctx context.Context, credentials *notesapi.Credentials) (*notesapi.User, error) {

	// Find user by username and password
	for _, user := range s.users {
		if user.Username == credentials.Username && user.Password == credentials.Password {
			return user, nil
		}
	}

	return nil, notesapi.ErrInvalidCredentials
}

// ListNotes returns all notes
func (s *InMemoryStore) ListNotes(ctx context.Context) ([]*notesapi.Note, error) {

	notes := make([]*notesapi.Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

// CreateNote creates a new note
func (s *InMemoryStore) CreateNote(ctx context.Context, note *notesapi.Note) (*notesapi.Note, error) {

	// Generate ID if not provided
	if note.Id == "" {
		note.Id = uuid.New().String()
	}

	// Store note
	s.notes[note.Id] = note
	return note, nil
}

// GetNote retrieves a note by ID
func (s *InMemoryStore) GetNote(ctx context.Context, id string) (*notesapi.Note, error) {

	note, exists := s.notes[id]
	if !exists {
		return nil, notesapi.ErrNoteNotFound
	}
	return note, nil
}

// UpdateNote updates a note
func (s *InMemoryStore) UpdateNote(ctx context.Context, id string, note *notesapi.Note) (*notesapi.Note, error) {

	_, exists := s.notes[id]
	if !exists {
		return nil, notesapi.ErrNoteNotFound
	}

	// Ensure ID is consistent
	note.Id = id
	s.notes[id] = note
	return note, nil
}

// DeleteNote deletes a note
func (s *InMemoryStore) DeleteNote(ctx context.Context, id string) error {

	_, exists := s.notes[id]
	if !exists {
		return notesapi.ErrNoteNotFound
	}

	delete(s.notes, id)
	return nil
}
