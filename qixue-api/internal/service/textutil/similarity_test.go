package textutil

import "testing"

func TestAverageSimilarityIdentical(t *testing.T) {
	got := AverageSimilarity("学习强国", "学习强国")
	if got < 99 {
		t.Fatalf("expected near 100, got %v", got)
	}
}

func TestAverageSimilarityEmpty(t *testing.T) {
	if AverageSimilarity("", "abc") != 0 {
		t.Fatal("expected 0")
	}
}
