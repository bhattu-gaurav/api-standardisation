package restapi

import (
	"api-standardisation/store"
	"api-standardisation/tsp-output/notesapi"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	store store.Store
}

func NewHTTPServer(s store.Store) *HTTPServer {
	return &HTTPServer{
		store: s,
	}
}
func (s *HTTPServer) AuthLogin(ctx echo.Context) error {
	var creds notesapi.Credentials

	if err := ctx.Bind(&creds); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	pbCreds := &notesapi.Credentials{
		Username: creds.Username,
		Password: creds.Password,
	}
	user, err := s.store.LoginUser(ctx.Request().Context(), pbCreds)
	if err != nil {
		if err == notesapi.ErrInvalidCredentials {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "login failed")
	}
	restUser := &notesapi.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	}
	return ctx.JSON(http.StatusOK, restUser)
}
func (s *HTTPServer) AuthRegister(c echo.Context) error {
	var user notesapi.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Convert to protobuf model
	pbUser := &notesapi.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	}

	createdUser, err := s.store.RegisterUser(c.Request().Context(), pbUser)
	if err != nil {
		if err == notesapi.ErrUserAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, "user already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to register user")
	}

	// Convert to REST model
	restUser := notesapi.User{
		Id:       createdUser.Id,
		Username: createdUser.Username,
		Password: createdUser.Password,
	}

	return c.JSON(http.StatusOK, toUserDTO(&restUser))
}

// NotesListNotes handles the GET /notes endpoint
func (s *HTTPServer) NotesListNotes(c echo.Context) error {
	pbNotes, err := s.store.ListNotes(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list notes")
	}

	// Convert to REST model
	var restNotes []notesapi.Note
	for _, pbNote := range pbNotes {
		restNotes = append(restNotes, notesapi.Note{
			Id:      pbNote.Id,
			Title:   pbNote.Title,
			Content: pbNote.Content,
			UserId:  pbNote.UserId,
		})
	}

	return c.JSON(http.StatusOK, NotesListResponse{Notes: restNotes})
}

// NotesCreateNote handles the POST /notes endpoint
func (s *HTTPServer) NotesCreateNote(c echo.Context) error {
	var note notesapi.Note
	if err := c.Bind(&note); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Convert to protobuf model
	pbNote := &notesapi.Note{
		Id:      note.Id,
		Title:   note.Title,
		Content: note.Content,
		UserId:  note.UserId,
	}

	createdNote, err := s.store.CreateNote(c.Request().Context(), pbNote)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create note")
	}

	// Convert to REST model
	restNote := notesapi.Note{
		Id:      createdNote.Id,
		Title:   createdNote.Title,
		Content: createdNote.Content,
		UserId:  createdNote.UserId,
	}

	return c.JSON(http.StatusOK, toNoteDTO(&restNote))
}

// NotesGetNote handles the GET /notes/{id} endpoint
func (s *HTTPServer) NotesGetNote(c echo.Context, id string) error {
	pbNote, err := s.store.GetNote(c.Request().Context(), id)
	if err != nil {
		if err == notesapi.ErrNoteNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "note not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get note")
	}

	// Convert to REST model
	restNote := notesapi.Note{
		Id:      pbNote.Id,
		Title:   pbNote.Title,
		Content: pbNote.Content,
		UserId:  pbNote.UserId,
	}

	return c.JSON(http.StatusOK, toNoteDTO(&restNote))
}

// NotesUpdateNote handles the PUT /notes/{id} endpoint
func (s *HTTPServer) NotesUpdateNote(c echo.Context, id string) error {
	var note notesapi.Note
	if err := c.Bind(&note); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Convert to protobuf model
	pbNote := &notesapi.Note{
		Id:      note.Id,
		Title:   note.Title,
		Content: note.Content,
		UserId:  note.UserId,
	}

	updatedNote, err := s.store.UpdateNote(c.Request().Context(), id, pbNote)
	if err != nil {
		if err == notesapi.ErrNoteNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "note not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update note")
	}

	// Convert to REST model
	restNote := notesapi.Note{
		Id:      updatedNote.Id,
		Title:   updatedNote.Title,
		Content: updatedNote.Content,
		UserId:  updatedNote.UserId,
	}

	return c.JSON(http.StatusOK, toNoteDTO(&restNote))
}

// NotesDeleteNote handles the DELETE /notes/{id} endpoint
func (s *HTTPServer) NotesDeleteNote(c echo.Context, id string) error {
	err := s.store.DeleteNote(c.Request().Context(), id)
	if err != nil {
		if err == notesapi.ErrNoteNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "note not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete note")
	}

	return c.NoContent(http.StatusOK)
}

// Define DTO (Data Transfer Object) types without protobuf internals
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Note struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserId  string `json:"userId"`
}

type NotesListResponse struct {
	Notes []notesapi.Note `json:"notes"`
}

// Convert protobuf User to DTO
func toUserDTO(pbUser *notesapi.User) User {
	return User{
		Id:       pbUser.Id,
		Username: pbUser.Username,
		Password: pbUser.Password,
	}
}

// Convert protobuf Note to DTO
func toNoteDTO(pbNote *notesapi.Note) Note {
	return Note{
		Id:      pbNote.Id,
		Title:   pbNote.Title,
		Content: pbNote.Content,
		UserId:  pbNote.UserId,
	}
}
