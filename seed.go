//go:build ignore

// Seeder para gvnotes. Puebla la base de datos con datos de prueba.
//
// Uso:
//
//	go run seed.go [--reset]
//
// Flags:
//
//	--reset   Elimina todos los notebooks y notas existentes antes de insertar
package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"

	db "gvnotes/db/generated"
	"gvnotes/internal/config"

	"github.com/google/uuid"
)

//go:embed db/migrations
var migrations embed.FS

func main() {
	reset := flag.Bool("reset", false, "Elimina notebooks y notas existentes antes de insertar")
	flag.Parse()

	sqlDB, err := config.OpenDB(migrations)
	if err != nil {
		log.Fatalf("abrir base de datos: %v", err)
	}
	defer sqlDB.Close()

	if *reset {
		if err := resetData(sqlDB); err != nil {
			log.Fatalf("reset: %v", err)
		}
		fmt.Println("✓ Datos existentes eliminados")
	}

	q := db.New(sqlDB)
	ctx := context.Background()

	if err := seed(ctx, q); err != nil {
		log.Fatalf("seed: %v", err)
	}

	fmt.Println("\n✓ Seed completado")
}

// resetData elimina todos los notebooks (cascade borra las notas) y notas raíz.
func resetData(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(`DELETE FROM notes WHERE notebook_id IS NULL`)
	if err != nil {
		return err
	}
	_, err = sqlDB.Exec(`DELETE FROM notebooks WHERE parent_id IS NULL`)
	return err
}

func seed(ctx context.Context, q db.Querier) error {
	// -------------------------------------------------------------------------
	// Nivel raíz: 20 notebooks + 15 notas sueltas = 35 items
	// -------------------------------------------------------------------------
	rootNotebooks, err := seedRootNotebooks(ctx, q, 20)
	if err != nil {
		return fmt.Errorf("notebooks raíz: %w", err)
	}

	if err := seedRootNotes(ctx, q, 15); err != nil {
		return fmt.Errorf("notas raíz: %w", err)
	}

	// -------------------------------------------------------------------------
	// Notebook "grande": 50 notas para probar desbordamiento de lista
	// -------------------------------------------------------------------------
	bigNB, err := createNotebook(ctx, q, nil, "📚 Notebook Grande (50 notas)", 100)
	if err != nil {
		return fmt.Errorf("notebook grande: %w", err)
	}
	if err := seedNotesInNotebook(ctx, q, bigNB.ID, 50, "Nota larga"); err != nil {
		return fmt.Errorf("notas notebook grande: %w", err)
	}
	fmt.Printf("  ✓ Notebook grande: %q → 50 notas\n", bigNB.Title)

	// -------------------------------------------------------------------------
	// Notebooks con sub-notebooks y notas anidadas (reutiliza los 3 primeros)
	// -------------------------------------------------------------------------
	nested := []struct {
		title  string
		subNBs int
		notes  int
	}{
		{"Trabajo", 5, 10},
		{"Personal", 4, 8},
		{"Proyectos", 6, 12},
	}

	for i, spec := range nested {
		if i >= len(rootNotebooks) {
			break
		}
		nb := rootNotebooks[i]
		if _, err := q.UpdateNotebookTitle(ctx, db.UpdateNotebookTitleParams{
			ID:    nb.ID,
			Title: spec.title,
		}); err != nil {
			return err
		}

		for j := 0; j < spec.subNBs; j++ {
			sub, err := createNotebook(ctx, q, &nb.ID, fmt.Sprintf("%s / Sub %d", spec.title, j+1), int64(j))
			if err != nil {
				return err
			}
			if err := seedNotesInNotebook(ctx, q, sub.ID, 3, "Nota"); err != nil {
				return err
			}
		}

		if err := seedNotesInNotebook(ctx, q, nb.ID, spec.notes, "Nota"); err != nil {
			return err
		}

		fmt.Printf("  ✓ %q: %d sub-notebooks · %d notas\n", spec.title, spec.subNBs, spec.notes)
	}

	total := len(rootNotebooks) + 1 // +1 notebook grande
	fmt.Printf("\nResumen raíz:\n")
	fmt.Printf("  Notebooks: %d\n", total)
	fmt.Printf("  Notas sueltas: 15\n")
	fmt.Printf("  Total items en raíz: %d\n", total+15)

	return nil
}

func seedRootNotebooks(ctx context.Context, q db.Querier, n int) ([]db.Notebook, error) {
	names := []string{
		"Trabajo", "Personal", "Proyectos", "Investigación", "Archivado",
		"En progreso", "Revisión", "Borradores", "Publicado", "Privado",
		"Compartido", "Favoritos", "Temporal", "Reciente", "Antiguo",
		"Crítico", "Opcional", "Experimental", "Estable", "Deprecado",
	}
	notebooks := make([]db.Notebook, 0, n)
	for i := 0; i < n; i++ {
		title := names[i%len(names)]
		if i >= len(names) {
			title = fmt.Sprintf("%s %d", title, i/len(names)+1)
		}
		nb, err := createNotebook(ctx, q, nil, title, int64(i))
		if err != nil {
			return nil, err
		}
		notebooks = append(notebooks, nb)
	}
	fmt.Printf("  ✓ %d notebooks en raíz\n", n)
	return notebooks, nil
}

func seedRootNotes(ctx context.Context, q db.Querier, n int) error {
	topics := []string{
		"Reunión de equipo", "Ideas del proyecto", "Lista de tareas", "Recordatorio",
		"Investigación inicial", "Feedback del cliente", "Plan de release", "Bug report",
		"Propuesta técnica", "Resumen semanal", "Notas de standup", "Brainstorm",
		"Decisiones del sprint", "Métricas del producto", "Retrospectiva",
	}
	for i := 0; i < n; i++ {
		title := topics[i%len(topics)]
		if i >= len(topics) {
			title = fmt.Sprintf("%s %d", title, i/len(topics)+1)
		}
		if _, err := createNote(ctx, q, nil, title, int64(i)); err != nil {
			return err
		}
	}
	fmt.Printf("  ✓ %d notas sueltas en raíz\n", n)
	return nil
}

func seedNotesInNotebook(ctx context.Context, q db.Querier, notebookID string, n int, prefix string) error {
	for i := 0; i < n; i++ {
		if _, err := createNote(ctx, q, &notebookID, fmt.Sprintf("%s %d", prefix, i+1), int64(i)); err != nil {
			return err
		}
	}
	return nil
}

func createNotebook(ctx context.Context, q db.Querier, parentID *string, title string, position int64) (db.Notebook, error) {
	return q.CreateNotebook(ctx, db.CreateNotebookParams{
		ID:       uuid.New().String(),
		ParentID: parentID,
		Title:    title,
		Position: position,
	})
}

func createNote(ctx context.Context, q db.Querier, notebookID *string, title string, position int64) (db.Note, error) {
	content := fmt.Sprintf(
		`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Contenido de ejemplo para: %s"}]}]}`,
		title,
	)
	return q.CreateNote(ctx, db.CreateNoteParams{
		ID:         uuid.New().String(),
		NotebookID: notebookID,
		Title:      title,
		Content:    content,
		Position:   position,
	})
}

