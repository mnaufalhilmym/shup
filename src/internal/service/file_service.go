package service

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"shup.hilmy.dev/src/internal/constant/errormsg"
	"shup.hilmy.dev/src/internal/entity"
	"shup.hilmy.dev/src/internal/repository"
)

type File struct {
	db   *sql.DB
	repo *repository.File
}

func NewFile(db *sql.DB, repo *repository.File) *File {
	return &File{db, repo}
}

func (s *File) Save(ctx context.Context, fileName string, body io.Reader) (*entity.File, error) {
	db, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg(errormsg.DATABASE_ERROR)
		return nil, err
	}
	defer db.Rollback()

	buf := make([]byte, 512)
	n, err := body.Read(buf)
	if err != nil && err != io.EOF {
		log.Error().Err(err).Msg("Failed to read the file bytes")
		return nil, err
	}

	// Detect content type from the buffer
	contentType := http.DetectContentType(buf[:n])

	// Reconstruct body so it don't lose the first bytes
	body = io.MultiReader(bytes.NewReader(buf[:n]), body)

	if err := os.MkdirAll("./files", os.ModePerm); err != nil {
		log.Error().Err(err).Msg("Failed to create the file directory")
	}

	fileID := uuid.New()
	dstFile, err := os.Create("./files/" + fileID.String())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create the file")
		return nil, err
	}
	defer dstFile.Close()

	size, err := io.Copy(dstFile, body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save the file bytes")
		return nil, err
	}

	createdAt := time.Now()

	fileEntity := &entity.File{
		ID:          fileID,
		FileName:    fileName,
		ContentType: contentType,
		Size:        size,
		CreatedAt:   createdAt,
		ExpiredAt:   createdAt.Add(24 * time.Hour),
	}

	if err := s.repo.Create(ctx, db, fileEntity); err != nil {
		log.Error().Err(err).Msg("Failed to save the file")
		return nil, err
	}

	if err := db.Commit(); err != nil {
		log.Error().Err(err).Msg(errormsg.DATABASE_ERROR)
		return nil, err
	}

	return fileEntity, nil
}

func (s *File) Get(ctx context.Context, fileID uuid.UUID) (*entity.File, error) {
	db, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg(errormsg.DATABASE_ERROR)
		return nil, err
	}
	defer db.Rollback()

	file, err := s.repo.FindByID(ctx, db, fileID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get the file")
		return nil, err
	}

	if err := db.Commit(); err != nil {
		log.Error().Err(err).Msg(errormsg.DATABASE_ERROR)
		return nil, err
	}

	return file, nil
}

func (s *File) DeleteAllExpired(ctx context.Context) error {
	db, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg(errormsg.DATABASE_ERROR)
		return err
	}
	defer db.Rollback()

	files, err := s.repo.FindAllExpired(ctx, db)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find all expired files")
		return err
	}

	deletedFileIDs := make([]uuid.UUID, 0, len(files))
	for _, file := range files {
		if err := os.Remove("./files/" + file.ID.String()); err != nil {
			return err
		}
		deletedFileIDs = append(deletedFileIDs, file.ID)
	}

	if err := s.repo.DeleteMany(ctx, db, deletedFileIDs); err != nil {
		return err
	}

	if err := db.Commit(); err != nil {
		log.Error().Err(err).Msg(errormsg.DATABASE_ERROR)
		return err
	}

	return nil
}
