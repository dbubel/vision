package store

import (
	"context"
	"github.com/dbubel/vision/models/requests"
	"github.com/dbubel/vision/models/tables"
)

func (s *Store) CreateWorkspace(ctx context.Context, cwr requests.CreateWorkspaceRequest) (*tables.Workspaces, error) {
	var w tables.Workspaces
	err := s.DbWriter.QueryRowxContext(ctx, "INSERT INTO workspaces (name, description) VALUES ($1, $2) RETURNING *", cwr.Name, cwr.Description).StructScan(&w)
	return &w, err
}
