package vectordb

import (
	"github.com/takara-ai/serverlessVector"
)

func CreateDB() *serverlessVector.VectorDB {
	db := serverlessVector.NewVectorDB(384)
	return db
}
