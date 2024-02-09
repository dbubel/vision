package tables

import "time"

type VectorRepo struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	WorkspaceID int       `db:"workspace_id"`
	Description string    `db:"description"`
	Dimension   int       `db:"dimension"`
	VectorTable string    `db:"vector_table"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   *NullTime `db:"updated_at"`
}
