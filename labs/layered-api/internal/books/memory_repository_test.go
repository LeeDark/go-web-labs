package books

import (
	"context"
	"sync"
	"testing"
)

func TestMemoryRepositoryListReturnsDeterministicSeed(t *testing.T) {
	repository := NewMemoryRepository()

	got, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	want := []Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan"},
		{ID: 2, Title: "Let's Go", Author: "Alex Edwards"},
	}
	if len(got) != len(want) {
		t.Fatalf("List() returned %d books, want %d", len(got), len(want))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("List()[%d] = %+v, want %+v", index, got[index], want[index])
		}
	}
}

func TestMemoryRepositoryCreateAssignsUniqueIDsConcurrently(t *testing.T) {
	const creates = 24

	repository := NewMemoryRepository()
	start := make(chan struct{})
	created := make(chan Book, creates)
	var waitGroup sync.WaitGroup

	for index := range creates {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			<-start
			book, err := repository.Create(context.Background(), CreateBookInput{
				Title:  "Concurrent book",
				Author: "Test author",
			})
			if err != nil {
				t.Errorf("Create() error = %v", err)
				return
			}
			created <- book
		}(index)
	}

	close(start)
	waitGroup.Wait()
	close(created)

	ids := make(map[int64]struct{}, creates)
	for book := range created {
		if book.ID < 3 || book.ID > creates+2 {
			t.Fatalf("created ID = %d, want range [3, %d]", book.ID, creates+2)
		}
		if _, exists := ids[book.ID]; exists {
			t.Fatalf("duplicate server-owned ID %d", book.ID)
		}
		ids[book.ID] = struct{}{}
	}
	if len(ids) != creates {
		t.Fatalf("created %d unique IDs, want %d", len(ids), creates)
	}

	listed, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got, want := len(listed), creates+2; got != want {
		t.Fatalf("List() returned %d books, want %d", got, want)
	}
	for index, book := range listed {
		if want := int64(index + 1); book.ID != want {
			t.Fatalf("List()[%d].ID = %d, want %d", index, book.ID, want)
		}
	}
}
