package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/dbubel/vision/models/requests"
	"github.com/dbubel/vision/models/tables"
	"github.com/pgvector/pgvector-go"
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
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (embedding vector(%d), meta JSONB NOT NULL DEFAULT '{}')`, tableName, dim)
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
		query := fmt.Sprintf("INSERT INTO %s (embedding,meta) VALUES ($1, $2)", tableName)
		rawJson, err := json.Marshal(i.Metadata)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err2 := tx.ExecContext(ctx, query, pgvector.NewVector(i.Data), rawJson)
		if err2 != nil {
			tx.Rollback()
			return err2
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = s.DbWriter.ExecContext(ctx, fmt.Sprintf("DROP INDEX IF EXISTS %s_embedding_idx", tableName))
	if err != nil {
		return err
	}

	n, err := s.GetIndexListSize(tableName)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("CREATE INDEX ON %s USING ivfflat (embedding vector_l2_ops) WITH (lists = %d);", tableName, n)
	_, err = s.DbWriter.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = s.DbWriter.ExecContext(ctx, fmt.Sprintf("VACUUM %s", tableName))
	if err != nil {
		return err
	}
	return nil
}

// GetIndexListSize is used for creating a vector index because the total number of rows matters
func (s *Store) GetIndexListSize(tableName string) (int, error) {
	var tableSize int
	err := s.DbReader.Get(&tableSize, fmt.Sprintf(`SELECT count(1) AS tableSize FROM %s`, tableName))
	if err != nil {
		return tableSize, err
	}

	if tableSize < 1000 {
		return 1, nil
	}

	return tableSize / 1000, nil
}

func getRandomTableName(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}

	return "vec_" + string(randomString)
}
