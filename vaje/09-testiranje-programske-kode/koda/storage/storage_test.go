package storage

import (
	"fmt"
	"testing"
)

func TestCreateAndReadAll(t *testing.T) {
	storage := NewTodoStorage()
	todo := Todo{Task: "task1", Completed: false}

	if err := storage.Create(&todo, &struct{}{}); err != nil {
		t.Fatalf("create returned error: %v", err)
	}

	got := make(map[string]Todo)
	if err := storage.Read(&Todo{}, &got); err != nil {
		t.Fatalf("read all returned error: %v", err)
	}

	if len(got) != 1 {
		t.Errorf("expected 1 todo, got %d", len(got))
	}

	if gotTodo, ok := got[todo.Task]; !ok {
		t.Errorf("expected task %q present", todo.Task)
	} else if gotTodo != todo {
		t.Errorf("read todo mismatch: got %+v want %+v", gotTodo, todo)
	}
}
func newStorageWithTodos(t *testing.T, todos ...Todo) *TodoStorage {
	t.Helper()
	storage := NewTodoStorage()
	for _, td := range todos {
		if err := storage.Create(&td, &struct{}{}); err != nil {
			t.Fatalf("create failed for %v: %v", td.Task, err)
		}
	}
	return storage
}
func TestUpdate(t *testing.T) {
	storage := newStorageWithTodos(t, Todo{Task: "present", Completed: false})

	t.Run("success", func(t *testing.T) {
		updated := Todo{Task: "present", Completed: true}
		if err := storage.Update(&updated, &struct{}{}); err != nil {
			t.Fatalf("update returned error: %v", err)
		}
		got := make(map[string]Todo)
		if err := storage.Read(&Todo{Task: "present"}, &got); err != nil {
			t.Fatalf("read returned error after update: %v", err)
		}
		if gotTodo := got["present"]; gotTodo.Completed != true {
			t.Fatalf("expected updated completion state, got %+v", gotTodo)
		}
	})
	t.Run("not found", func(t *testing.T) {
		err := storage.Update(&Todo{Task: "missing", Completed: true}, &struct{}{})
		if err != ErrorNotFound {
			t.Fatalf("expected error for missing task update, got %v", err)
		}
	})
}

func FuzzCreateReadDeleteTest(f *testing.F) {
	// seme za generiranje vhodnih podatkov
	f.Add("seed", false)
	f.Fuzz(func(t *testing.T, task string, completed bool) {
		storage := NewTodoStorage()
		todo := Todo{Task: task, Completed: completed}
		t.Logf("Testing with %+v", todo)
		if err := storage.Create(&todo, &struct{}{}); err != nil {
			t.Fatalf("create failed: %v", err)
		}
		got := make(map[string]Todo)
		if err := storage.Read(&Todo{Task: task}, &got); err != nil {
			t.Fatalf("read failed: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected one entry, got %d", len(got))
		}
		if got[task].Completed != completed {
			t.Fatalf("completion mismatch: got %v want %v", got[task].Completed, completed)
		}
		if err := storage.Delete(&Todo{Task: task}, &struct{}{}); err != nil {
			t.Fatalf("delete failed: %v", err)
		}
		if err := storage.Read(&Todo{Task: task}, &got); err != ErrorNotFound {
			t.Fatalf("expected not found after delete, got %v", err)
		}
	})
}

func BenchmarkCreate(b *testing.B) {
	storage := NewTodoStorage()
	i := 0
	b.ResetTimer()
	for b.Loop() {
		i++
		if err := storage.Create(&Todo{Task: fmt.Sprintf("task-%d", i)}, &struct{}{}); err != nil {
			b.Fatalf("create failed: %v", err)
		}
	}
}
