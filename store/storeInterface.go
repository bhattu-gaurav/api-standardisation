package store

import (
	"api-standardisation/tsp-output/notesapi"
	"context"
)

type Store interface {
	RegisterUser(ctx context.Context, user *notesapi.User) (*notesapi.User, error)
	LoginUser(ctx context.Context, credentials *notesapi.Credentials) (*notesapi.User, error)

	// Note methods
	ListNotes(ctx context.Context) ([]*notesapi.Note, error)
	CreateNote(ctx context.Context, note *notesapi.Note) (*notesapi.Note, error)
	GetNote(ctx context.Context, id string) (*notesapi.Note, error)
	UpdateNote(ctx context.Context, id string, note *notesapi.Note) (*notesapi.Note, error)
	DeleteNote(ctx context.Context, id string) error
}
