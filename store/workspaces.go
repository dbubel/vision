package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/dbubel/vision/models/requests"
	"github.com/dbubel/vision/models/tables"
	"github.com/pgvector/pgvector-go"
)

func (s *Store) CreateWorkspace(ctx context.Context, cwr requests.CreateWorkspaceRequest) (*tables.Workspaces, error) {
	var w tables.Workspaces
	err := s.DbWriter.QueryRowxContext(ctx, "INSERT INTO public.workspaces (name,  description) VALUES ($1,  $2) RETURNING *", cwr.Name, cwr.Description).StructScan(&w)
	return &w, err
}

func (s *Store) CreateVectorTableMapping(ctx context.Context, cwr requests.CreateRepoRequest) (*tables.VectorRepo, error) {
	var w tables.VectorRepo
	tableName := generateRandomString(25)
	err := s.DbWriter.QueryRowxContext(ctx, `
    INSERT INTO workspace_vector_mappings ( workspace_id, dimension, vector_table,name,updated_at)
    SELECT id, $2,$3 ,$4,NOW()
    FROM workspaces
    WHERE name = $1 RETURNING *`, cwr.Workspace, cwr.Dimension, tableName, cwr.Name).StructScan(&w)
	// if err != nil {
	//   return &w, err
	// }

	// err = s.createVectorTable(ctx, tableName, cwr.Dimension)
	// if err != nil {
	//    return &w, err
	// }
	return &w, err
}

func (s *Store) CreateVectorTable(ctx context.Context, tableName string, dim int) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (embedding vector(%d))`, tableName, dim)
	_, err := s.DbWriter.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return err
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}

	return "vec_" + string(randomString)
}

func (s *Store) ValidateInsert(ctx context.Context, workspace, repoName string) (string, error) {
	var wks string
	return wks, s.DbReader.QueryRowx(
		`select wvm.vector_table from workspaces w
join workspace_vector_mappings wvm on w.id = wvm.workspace_id
where w.name = $1 and wvm.name = $2;`, workspace, repoName).Scan(&wks)
}

func (s *Store) InsertVectors(ctx context.Context, tableName string, vectors requests.Vectors) error {
	// pgvec col requires float32
	// var embeddingsFloat32 []float32
	// embeddingsFloat32 = make([]float32, len(content.Embeddings))
	//
	// for i := 0; i < len(content.Embeddings); i++ {
	// 	embeddingsFloat32[i] = float32(content.Embeddings[i])
	// }
	//
	for _, i := range vectors.Data {
		query := fmt.Sprintf("INSERT INTO %s ( embedding) VALUES($1) ", tableName)
		_, err := s.DbWriter.ExecContext(ctx, query, pgvector.NewVector(i))
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	return nil
}
