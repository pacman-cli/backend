package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/puspo/basicnewproject/internal/models"
)

// PostRepository exposes CRUD methods for posts.
type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create inserts a new post and returns its ID.
func (r *PostRepository) Create(ctx context.Context, p *models.Post) (int64, error) {
	// Marshal JSON fields
	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return 0, err
	}
	metadataJSON, err := json.Marshal(p.Metadata)
	if err != nil {
		return 0, err
	}

	// Use context with timeout to avoid hanging queries.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO posts (title, content, tags, metadata, created_at, updated_at)
		 VALUES (?, ?, CAST(? AS JSON), CAST(? AS JSON), NOW(), NOW())`,
		p.Title, p.Content, string(tagsJSON), string(metadataJSON),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Get retrieves a post by id.
func (r *PostRepository) Get(ctx context.Context, id int64) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `SELECT id, title, content, tags, metadata, created_at, updated_at FROM posts WHERE id = ?`, id)

	var (
		p           models.Post
		tagsRaw     sql.NullString
		metadataRaw sql.NullString
	)
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &tagsRaw, &metadataRaw, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if tagsRaw.Valid {
		_ = json.Unmarshal([]byte(tagsRaw.String), &p.Tags)
	}
	if metadataRaw.Valid {
		_ = json.Unmarshal([]byte(metadataRaw.String), &p.Metadata)
	}
	return &p, nil
}

// List returns posts with basic pagination: offset/limit.
func (r *PostRepository) List(ctx context.Context, offset, limit int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `SELECT id, title, content, tags, metadata, created_at, updated_at FROM posts ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var (
			p           models.Post
			tagsRaw     sql.NullString
			metadataRaw sql.NullString
		)
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &tagsRaw, &metadataRaw, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if tagsRaw.Valid {
			_ = json.Unmarshal([]byte(tagsRaw.String), &p.Tags)
		}
		if metadataRaw.Valid {
			_ = json.Unmarshal([]byte(metadataRaw.String), &p.Metadata)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// Update modifies an existing post. Returns false if the row doesn't exist.
func (r *PostRepository) Update(ctx context.Context, id int64, p *models.Post) (bool, error) {
	// Marshal JSON fields
	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return false, err
	}
	metadataJSON, err := json.Marshal(p.Metadata)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(
		ctx,
		`UPDATE posts SET title = ?, content = ?, tags = CAST(? AS JSON), metadata = CAST(? AS JSON), updated_at = NOW() WHERE id = ?`,
		p.Title, p.Content, string(tagsJSON), string(metadataJSON), id,
	)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

// Delete removes a post by id. Returns false if the row doesn't exist.
func (r *PostRepository) Delete(ctx context.Context, id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
