package tables

import (
	"encoding/json"

	"github.com/pgvector/pgvector-go"
)

type Vector struct {
	Embedding pgvector.Vector `json:"embedding" db:"embedding"`
	Metadata  json.RawMessage `json:"metadata" db:"meta"`
}

type Vectors []Vector

// this is dumb but for some reason the pgvector.Vector
// does not marshal correctly on responses. so we have to
// call slice() and build a new struct.
func (v *Vector) MarshalJSON() ([]byte, error) {
	x := struct {
		Embedding []float32
		Metadata  json.RawMessage
	}{
		Embedding: v.Embedding.Slice(),
		Metadata:  v.Metadata,
	}
	return json.Marshal(x)
}
