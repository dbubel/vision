package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dbubel/vision/models/requests"
	"github.com/dbubel/vision/models/tables"
	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
	"math/rand"
	"time"
)

func (s *Store) CreateVectorTableMapping(ctx context.Context, crr requests.CreateRepoRequest) (*tables.VectorRepo, error) {
	var w tables.VectorRepo
	tableName := getRandomTableName(25)
	err := s.DbWriter.QueryRowxContext(ctx, `
		INSERT INTO workspace_vector_mappings (workspace_id, dimension, vector_table, name, updated_at)
		SELECT id, $2, $3, $4, NOW()
		FROM workspaces
		WHERE name = $1 RETURNING *`, crr.Workspace, crr.Dimension, tableName, crr.Name).StructScan(&w)
	return &w, err
}

func (s *Store) CreateVectorTable(ctx context.Context, tableName string, dim int) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (embedding vector(%d))`, tableName, dim)
	_, err := s.DbWriter.ExecContext(ctx, query)
	return err
}

func (s *Store) ValidateInsert(ctx context.Context, workspace, repoName string) (string, error) {
	var wks string
	return wks, s.DbReader.QueryRowx(
		`SELECT wvm.vector_table FROM workspaces w
		JOIN workspace_vector_mappings wvm ON w.id = wvm.workspace_id
		WHERE w.name = $1 AND wvm.name = $2;`, workspace, repoName).Scan(&wks)
}

func (s *Store) InsertVectors(ctx context.Context, tableName string, vectors requests.Vectors) error {
	// Do this potentially large insert in a transaction - I heard it's faster this way.
	tx, err := s.DbWriter.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	for _, i := range vectors.Data {
		query := fmt.Sprintf("INSERT INTO %s (embedding) VALUES ($1)", tableName)
		_, err := tx.ExecContext(ctx, query, pgvector.NewVector(i))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.ExecContext(ctx, fmt.Sprintf("DROP INDEX %s_embedding_idx", tableName))
	if err != nil {
		tx.Rollback()
		return err
	}

	n, err := s.GetIndexListSize(tx, tableName)
	if err != nil {
		tx.Rollback()
		return err
	}
	query := fmt.Sprintf("CREATE INDEX ON %s USING hnsw (embedding vector_l2_ops) WITH (lists = %d);", tableName, n)
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, fmt.Sprintf("VACUUM %s", tableName))
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// GetIndexListSize is used for creating a vector index because the total number of rows matters
func (s *Store) GetIndexListSize(tx *sqlx.Tx, tableName string) (int, error) {
	var tableSize int
	err := tx.Get(&tableSize, fmt.Sprintf(`SELECT count(1) AS tableSize FROM %s`, tableName))
	if err != nil {
		return tableSize, err
	}

	if tableSize < 1000 {
		return 1, nil
	}

	return tableSize / 1000, nil
}

func getRandomTableName(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}

	return "vec_" + string(randomString)
}
