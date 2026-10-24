package menutree

import (
	"testing"

	"qixue-api/internal/model"
)

func intPtr(n int) *int { return &n }

func TestTreeRootsAndChildren(t *testing.T) {
	rows := []model.AdminMenu{
		{ID: 2, Name: "子", ParentID: intPtr(1)},
		{ID: 1, Name: "根"},
		{ID: 3, Name: "另一根", ParentID: intPtr(0)},
		{ID: 9, Name: "孤儿", ParentID: intPtr(99)},
	}
	got := Tree(rows)
	if len(got) != 2 {
		t.Fatalf("roots=%d, want 2", len(got))
	}
	if got[0].ID != 1 || len(got[0].Children) != 1 || got[0].Children[0].ID != 2 {
		t.Fatalf("first root = %+v", got[0])
	}
	if got[1].ID != 3 {
		t.Fatalf("second root = %d", got[1].ID)
	}
}
