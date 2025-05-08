package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"shup.hilmy.dev/src/internal/entity"
)

type File struct{}

func NewFile(db *sql.DB) *File {
	f := &File{}
	f.createTable(db)
	return f
}

func (*File) Create(ctx context.Context, db *sql.Tx, file *entity.File) error {
	query := "INSERT INTO files (" +
		"id, file_name, content_type, size, created_at, expired_at" +
		") VALUES (" +
		"$1, $2, $3, $4, $5, $6" +
		")"

	if _, err := db.ExecContext(
		ctx,
		query,
		file.ID,
		file.FileName,
		file.ContentType,
		file.Size,
		file.CreatedAt,
		file.ExpiredAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *File) FindByID(ctx context.Context, db *sql.Tx, id uuid.UUID) (*entity.File, error) {
	query := "SELECT id, file_name, content_type, size, created_at, expired_at" +
		" FROM files" +
		" WHERE id = $1"

	row := db.QueryRowContext(ctx, query, id)

	return r.getEntity(row)
}

func (r *File) FindAllExpired(ctx context.Context, db *sql.Tx) ([]entity.File, error) {
	query := "SELECT id, file_name, content_type, size, created_at, expired_at" +
		" FROM files" +
		" WHERE expired_at < CURRENT_TIMESTAMP" +
		" ORDER BY expired_at DESC"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.getEntities(rows)
}

func (*File) DeleteMany(ctx context.Context, db *sql.Tx, ids []uuid.UUID) error {
	query := "DELETE FROM files" +
		" WHERE id = ANY($1)"

	if _, err := db.ExecContext(ctx, query, ids); err != nil {
		return err
	}

	return nil
}

func (*File) createTable(db *sql.DB) {
	query := "CREATE TABLE IF NOT EXISTS files (" +
		" id TEXT PRIMARY KEY," +
		" file_name TEXT NOT NULL," +
		" content_type TEXT NOT NULL," +
		" size INTEGER NOT NULL," +
		" created_at DATETIME NOT NULL," +
		" expired_at DATETIME NOT NULL" +
		")"

	if _, err := db.Exec(query); err != nil {
		log.Fatal().Err(err).Msg("Error initialize database")
	}
}

func (*File) getEntity(row *sql.Row) (*entity.File, error) {
	var file entity.File

	if err := row.Scan(
		&file.ID,
		&file.FileName,
		&file.ContentType,
		&file.Size,
		&file.CreatedAt,
		&file.ExpiredAt,
	); err != nil {
		return nil, err
	}

	return &file, nil
}

func (*File) getEntities(rows *sql.Rows) ([]entity.File, error) {
	var files []entity.File

	for rows.Next() {
		var file entity.File
		if err := rows.Scan(
			&file.ID,
			&file.FileName,
			&file.ContentType,
			&file.Size,
			&file.CreatedAt,
			&file.ExpiredAt,
		); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
