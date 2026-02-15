package metadata

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kenny/lessons/demo/pb"

	"github.com/google/uuid"
)

type LessonStatus string

const (
	StatusDraft     LessonStatus = "DRAFT"
	StatusPublished LessonStatus = "PUBLISHED"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateLesson(ctx context.Context, contentID string) (rootID string, versionNumber int32, err error) {
	rootID = uuid.NewString()
	versionID := uuid.NewString()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", 0, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO lesson_roots (id, created_at) VALUES ($1, $2)`,
		rootID, time.Now(),
	)
	if err != nil {
		return "", 0, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO lesson_versions (id, root_id, version_number, status, content_id, created_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		versionID, rootID, 1, string(StatusDraft), contentID, time.Now(),
	)
	if err != nil {
		return "", 0, err
	}

	if err := tx.Commit(); err != nil {
		return "", 0, err
	}

	return rootID, 1, nil
}

func (s *Store) CreateVersion(ctx context.Context, rootID, contentID string) (int32, error) {
	var latest int32
	err := s.db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(version_number), 0) FROM lesson_versions WHERE root_id = $1`,
		rootID,
	).Scan(&latest)
	if err != nil {
		return 0, err
	}

	versionID := uuid.NewString()
	versionNumber := latest + 1

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO lesson_versions (id, root_id, version_number, status, content_id, created_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		versionID, rootID, versionNumber, string(StatusDraft), contentID, time.Now(),
	)
	if err != nil {
		return 0, err
	}

	return versionNumber, nil
}

func (s *Store) PublishVersion(ctx context.Context, rootID string, versionNumber int32) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE lesson_versions
         SET status = $1, published_at = $2
         WHERE root_id = $3 AND version_number = $4 AND status != $1`,
		string(StatusPublished), time.Now(), rootID, versionNumber,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrVersionNotFound
	}
	return nil
}

func (s *Store) GetVersion(ctx context.Context, rootID string, versionNumber int32) (*pb.GetVersionResponse, error) {
	var contentID string
	err := s.db.QueryRowContext(ctx,
		`SELECT content_id FROM lesson_versions WHERE root_id = $1 AND version_number = $2`,
		rootID, versionNumber,
	).Scan(&contentID)
	if err == sql.ErrNoRows {
		return nil, ErrVersionNotFound
	}
	if err != nil {
		return nil, err
	}

	return &pb.GetVersionResponse{
		RootId:        rootID,
		VersionNumber: versionNumber,
		ContentId:     contentID,
	}, nil
}
