package vectordb

import "testing"

func TestCreateVectorDB(t *testing.T) {
	db := CreateDB()
	if db.Size() > 0 {
		t.Errorf("Expecting new DB to have size zero, got %d", db.Size())
	}
	if dims, ok := db.GetStats()["dimension"]; ok {
		if dimsInt, isInt := dims.(int); isInt {
			if dimsInt != 384 {
				t.Errorf("Expecting 'dimension' to be 384, but got %d", dimsInt)
			}
		} else {
			t.Error("Expecting 'dimension' to be an integer, but they are not.")
		}
	} else {
		t.Error("Expecting 'dimension' to be found in stats, but it is not.")
	}
}
