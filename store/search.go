package store

import (
	"context"
	"fmt"

	"github.com/dbubel/vision/models/tables"
	"github.com/pgvector/pgvector-go"
)

func (s *Store) Search(ctx context.Context, tableName string, vector []float32) (*tables.Vector, error) {
	var w tables.Vector
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY embedding <-> $1 LIMIT 1", tableName)
	err := s.DbWriter.GetContext(ctx, &w, query, pgvector.NewVector(vector))
	return &w, err
}
