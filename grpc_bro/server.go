package grpc_bro

import (
	"api-standardisation/store"
	"api-standardisation/tsp-output/notesapi"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GRPCServer implements the gRPC service interfaces
type GRPCServer struct {
	notesapi.UnimplementedAuthServer
	notesapi.UnimplementedNotesServer
	store store.Store
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(store store.Store) *GRPCServer {
	return &GRPCServer{
		store: store,
	}
}

// Register implements the Register RPC method
func (s *GRPCServer) Register(ctx context.Context, req *notesapi.RegisterRequest) (*notesapi.User, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user is required")
	}

	user, err := s.store.RegisterUser(ctx, req.User)
	if err != nil {
		if err == notesapi.ErrUserAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return user, nil
}

// Login implements the Login RPC method
func (s *GRPCServer) Login(ctx context.Context, req *notesapi.LoginRequest) (*notesapi.User, error) {
	if req.Credentials == nil {
		return nil, status.Errorf(codes.InvalidArgument, "credentials are required")
	}

	user, err := s.store.LoginUser(ctx, req.Credentials)
	if err != nil {
		if err == notesapi.ErrInvalidCredentials {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return user, nil
}

// ListNotes implements the ListNotes RPC method
func (s *GRPCServer) ListNotes(ctx context.Context, _ *emptypb.Empty) (*notesapi.ListNotesResponse, error) {
	notes, err := s.store.ListNotes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	convertedNotes := make([]*notesapi.Note, len(notes))
	for i, note := range notes {
		convertedNotes[i] = &notesapi.Note{
			Id:      note.Id,
			Title:   note.Title,
			Content: note.Content,
		}
	}
	return &notesapi.ListNotesResponse{Notes: convertedNotes}, nil
}

// CreateNote implements the CreateNote RPC method
func (s *GRPCServer) CreateNote(ctx context.Context, req *notesapi.CreateNoteRequest) (*notesapi.Note, error) {
	if req.Note == nil {
		return nil, status.Errorf(codes.InvalidArgument, "note is required")
	}

	note, err := s.store.CreateNote(ctx, req.Note)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create note: %v", err)
	}

	return note, nil
}

// GetNote implements the GetNote RPC method
func (s *GRPCServer) GetNote(ctx context.Context, req *notesapi.GetNoteRequest) (*notesapi.Note, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	note, err := s.store.GetNote(ctx, req.Id)
	if err != nil {
		if err == notesapi.ErrNoteNotFound {
			return nil, status.Errorf(codes.NotFound, "note not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get note: %v", err)
	}

	return note, nil
}

// UpdateNote implements the UpdateNote RPC method
func (s *GRPCServer) UpdateNote(ctx context.Context, req *notesapi.UpdateNoteRequest) (*notesapi.Note, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.Note == nil {
		return nil, status.Errorf(codes.InvalidArgument, "note is required")
	}

	note, err := s.store.UpdateNote(ctx, req.Id, req.Note)
	if err != nil {
		if err == notesapi.ErrNoteNotFound {
			return nil, status.Errorf(codes.NotFound, "note not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update note: %v", err)
	}

	return note, nil
}

// DeleteNote implements the DeleteNote RPC method
func (s *GRPCServer) DeleteNote(ctx context.Context, req *notesapi.DeleteNoteRequest) (*notesapi.Empty, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	err := s.store.DeleteNote(ctx, req.Id)
	if err != nil {
		if err == notesapi.ErrNoteNotFound {
			return nil, status.Errorf(codes.NotFound, "note not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete note: %v", err)
	}

	return &notesapi.Empty{}, nil
}
