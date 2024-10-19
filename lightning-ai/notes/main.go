package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Note struct {
	ID        string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DataAPI interface {
	ListNotes() ([]Note, error)
	GetNote(id string) (*Note, error)
	CreateNote(content string) (*Note, error)
	UpdateNote(id, content string) error
	DeleteNote(id string) error
}

type NoteService struct {
	db *sql.DB
}

var _ DataAPI = &NoteService{}

func NewNoteService(db *sql.DB) *NoteService {
	return &NoteService{db: db}
}

func (s *NoteService) ListNotes() ([]Note, error) {
	rows, err := s.db.Query("SELECT id, content, created_at, updated_at FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		err := rows.Scan(&n.ID, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (s *NoteService) GetNote(id string) (*Note, error) {
	var n Note
	err := s.db.QueryRow("SELECT id, content, created_at, updated_at FROM notes WHERE id = $1", id).
		Scan(&n.ID, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("note not found")
		}
		return nil, err
	}
	return &n, nil
}

func (s *NoteService) CreateNote(content string) (*Note, error) {
	n := Note{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := s.db.Exec("INSERT INTO notes (id, content, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		n.ID, n.Content, n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *NoteService) UpdateNote(id, content string) error {
	result, err := s.db.Exec("UPDATE notes SET content = $1, updated_at = $2 WHERE id = $3",
		content, time.Now(), id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("note not found")
	}
	return nil
}

func (s *NoteService) DeleteNote(id string) error {
	result, err := s.db.Exec("DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("note not found")
	}
	return nil
}

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", "postgres://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create the notes table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id VARCHAR(36) PRIMARY KEY,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create a new NoteService
	noteService := NewNoteService(db)

	// Create a new note
	newNote, err := noteService.CreateNote("This is a test note")
	if err != nil {
		log.Fatalf("Failed to create note: %v", err)
	}
	fmt.Printf("Created note: %+v\n", newNote)

	// Get the created note
	retrievedNote, err := noteService.GetNote(newNote.ID)
	if err != nil {
		log.Fatalf("Failed to get note: %v", err)
	}
	fmt.Printf("Retrieved note: %+v\n", retrievedNote)

	// Update the note
	err = noteService.UpdateNote(newNote.ID, "This is an updated test note")
	if err != nil {
		log.Fatalf("Failed to update note: %v", err)
	}
	fmt.Println("Note updated successfully")

	// List all notes
	notes, err := noteService.ListNotes()
	if err != nil {
		log.Fatalf("Failed to list notes: %v", err)
	}
	fmt.Println("All notes:")
	for _, note := range notes {
		fmt.Printf("%+v\n", note)
	}

	// Delete the note
	err = noteService.DeleteNote(newNote.ID)
	if err != nil {
		log.Fatalf("Failed to delete note: %v", err)
	}
	fmt.Println("Note deleted successfully")

	// Try to get the deleted note (should fail)
	_, err = noteService.GetNote(newNote.ID)
	if err != nil {
		fmt.Printf("Expected error when getting deleted note: %v\n", err)
	}
}