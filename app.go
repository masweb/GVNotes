package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"

	db "gvnotes/db/generated"
	"gvnotes/internal/config"
	"gvnotes/internal/controllers"
	"gvnotes/internal/dto"
	"gvnotes/internal/services"
)

// appMigrations is set from main.go before NewApp() is called.
var appMigrations embed.FS

// App is the Wails application struct. It acts as a facade over internal
// controllers, exposing methods to the frontend via Wails bindings.
type App struct {
	ctx            context.Context
	sqlDB          *sql.DB
	auth           *controllers.AuthController
	notebook       *controllers.NotebookController
	note           *controllers.NoteController
	image          *controllers.ImageController
	imageListener  net.Listener
	imageServerURL string
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. Dependencies are composed here.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	sqlDB, err := config.OpenDB(appMigrations)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	a.sqlDB = sqlDB

	imagesDir, err := config.ImagesDir()
	if err != nil {
		log.Fatalf("failed to resolve images directory: %v", err)
	}

	querier := db.New(sqlDB)

	a.auth = controllers.NewAuthController(services.NewAuthService(querier))
	a.notebook = controllers.NewNotebookController(services.NewNotebookService(querier))
	a.note = controllers.NewNoteController(services.NewNoteService(querier))
	a.image = controllers.NewImageController(services.NewImageService(querier, imagesDir))

	// Start a local HTTP server to serve image files. This is necessary because
	// the Wails AssetServer Handler does not work in dev mode with Vite v5+.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("failed to start image server: %v", err)
	}
	a.imageListener = ln
	a.imageServerURL = fmt.Sprintf("http://127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port)
	go func() {
		if err := http.Serve(ln, &imageFileHandler{imagesDir: imagesDir}); err != nil && err != http.ErrServerClosed {
			log.Printf("image server stopped: %v", err)
		}
	}()
}

// shutdown is called when the app closes.
func (a *App) shutdown(_ context.Context) {
	if a.imageListener != nil {
		a.imageListener.Close()
	}
	if a.sqlDB != nil {
		a.sqlDB.Close()
	}
}

// GetImageServerURL returns the base URL of the local image HTTP server.
func (a *App) GetImageServerURL() string {
	return a.imageServerURL
}

// --- Auth ---

func (a *App) GetAuthStatus() interface{} {
	return a.auth.GetAuthStatus()
}

func (a *App) SetPassword(password string) interface{} {
	return a.auth.SetPassword(password)
}

func (a *App) VerifyPassword(password string) interface{} {
	return a.auth.VerifyPassword(password)
}

// --- Notebooks ---

func (a *App) ListNotebooks(parentID string) ([]dto.NotebookListItem, error) {
	return a.notebook.ListNotebooks(parentID)
}

func (a *App) GetNotebook(id string) (dto.NotebookDetail, error) {
	return a.notebook.GetNotebook(id)
}

func (a *App) CreateNotebook(req dto.CreateNotebookRequest) (dto.NotebookDetail, error) {
	return a.notebook.CreateNotebook(req)
}

func (a *App) UpdateNotebookTitle(id string, req dto.UpdateNotebookTitleRequest) (dto.NotebookDetail, error) {
	return a.notebook.UpdateNotebookTitle(id, req)
}

func (a *App) UpdateNotebookPosition(id string, req dto.UpdatePositionRequest) error {
	return a.notebook.UpdateNotebookPosition(id, req)
}

func (a *App) DeleteNotebook(id string) error {
	return a.notebook.DeleteNotebook(id)
}

// --- Notes ---

func (a *App) ListNotes(notebookID string) ([]dto.NoteListItem, error) {
	return a.note.ListNotes(notebookID)
}

func (a *App) GetNote(id string) (dto.NoteDetail, error) {
	return a.note.GetNote(id)
}

func (a *App) CreateNote(req dto.CreateNoteRequest) (dto.NoteDetail, error) {
	return a.note.CreateNote(req)
}

func (a *App) UpdateNoteTitle(id string, req dto.UpdateNoteTitleRequest) (dto.NoteDetail, error) {
	return a.note.UpdateNoteTitle(id, req)
}

func (a *App) UpdateNoteContent(id string, req dto.UpdateNoteContentRequest) (dto.NoteDetail, error) {
	return a.note.UpdateNoteContent(id, req)
}

func (a *App) UpdateNotePosition(id string, req dto.UpdatePositionRequest) error {
	return a.note.UpdateNotePosition(id, req)
}

func (a *App) DeleteNote(id string) error {
	return a.note.DeleteNote(id)
}

// --- Images ---

func (a *App) ListImagesByNote(noteID string) ([]dto.ImageItem, error) {
	return a.image.ListImagesByNote(noteID)
}

func (a *App) GetImage(id string) (dto.ImageItem, error) {
	return a.image.GetImage(id)
}

func (a *App) SaveImage(noteID string, mimeType string, data []byte) (dto.ImageItem, error) {
	return a.image.SaveImage(noteID, mimeType, data)
}

func (a *App) DeleteImage(id string) error {
	return a.image.DeleteImage(id)
}

func (a *App) GetImagePath(filename string) (string, error) {
	return a.image.GetImagePath(filename)
}
