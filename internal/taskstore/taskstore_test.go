package taskstore

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T) {
	ts := New()
	mustParseDate := func(tstr string) time.Time {
		timeFormat := "2006-Jan-02"
		tt, err := time.Parse(timeFormat, tstr)
		if err != nil {
			t.Fatal(err)
		}
		return tt
	}

	firstId := ts.CreateTask("First task", []string{"first", "test"}, mustParseDate("2019-Jan-01"))

	// Test GetTask
	task, err := ts.GetTask(firstId)
	if err != nil {
		t.Fatal(err)
	}

	if task.Id != firstId {
		t.Errorf("GetTask(%d) returned a task with id %d", firstId, task.Id)
	}

	if task.Text != "First task" {
		t.Errorf("GetTask(%d) returned a task with text %q", firstId, task.Text)
	}

	_, err = ts.GetTask(100)
	if err == nil {
		t.Fatal("err should not be nil")
	}

	// Test GetAllTasks
	allTasks := ts.GetAllTasks()
	if len(allTasks) != 1 || allTasks[0].Id != firstId {
		t.Errorf("GetAllTasks() returned %d tasks, expected 1", len(allTasks))
	}

	// Test CreateTask and GetAllTasks
	ts.CreateTask("Second task", []string{"second", "test", "even"}, mustParseDate("2019-Jan-02"))
	ts.CreateTask("Third task", []string{"third", "test"}, mustParseDate("2019-Jan-02"))
	ts.CreateTask("Fourth task", []string{"fourth", "test", "even"}, mustParseDate("2019-Jan-03"))
	allTasks = ts.GetAllTasks()
	if len(allTasks) != 4 {
		t.Errorf("GetAllTasks() returned %d tasks, expected 4", len(allTasks))
	}

	// Test GetTasksByTag
	tasksByTag := ts.GetTasksByTag("test")
	if len(tasksByTag) != 4 {
		t.Errorf("GetTasksByTag(\"test\") returned %d tasks, expected 4", len(tasksByTag))
	}

	tasksByTag = ts.GetTasksByTag("even")
	if len(tasksByTag) != 2 {
		t.Errorf("GetTasksByTag(\"even\") returned %d tasks, expected 2", len(tasksByTag))
	}

	// Test GetTasksByDueDate
	tasksByDueDate := ts.GetTasksByDate(mustParseDate("2019-Jan-01").Date())
	if len(tasksByDueDate) != 1 {
		t.Errorf("GetTasksByDueDate(\"2019-Jan-01\") returned %d tasks, expected 1", len(tasksByDueDate))
	}

	tasksByDueDate = ts.GetTasksByDate(mustParseDate("2019-Jan-02").Date())
	if len(tasksByDueDate) != 2 {
		t.Errorf("GetTasksByDueDate(\"2019-Jan-02\") returned %d tasks, expected 2", len(tasksByDueDate))
	}
}

func TestCreateAndDelete(t *testing.T) {
	ts := New()

	firstId := ts.CreateTask("First task", []string{"first", "test"}, time.Now())
	ts.CreateTask("Second task", []string{"second", "test"}, time.Now())
	ts.CreateTask("Third task", []string{"third", "test"}, time.Now())
	ts.CreateTask("Fourth task", []string{"fourth", "test"}, time.Now())

	// Test DeleteTask
	if err := ts.DeleteTask(1000); err == nil {
		t.Fatalf("DeleteTask(1000) should return an error")
	}

	if err := ts.DeleteTask(firstId); err != nil {
		t.Fatal(err)
	}

	if err := ts.DeleteTask(firstId); err == nil {
		t.Fatalf("DeleteTask(%d) should return an error", firstId)
	}

	allTasks := ts.GetAllTasks()
	if len(allTasks) != 3 {
		t.Fatalf("GetAllTasks() returned %d tasks, expected 3", len(allTasks))
	}

	// Test DeleteAllTasks
	ts.DeleteAllTasks()
	allTasks = ts.GetAllTasks()
	if len(allTasks) != 0 {
		t.Fatalf("GetAllTasks() returned %d tasks, expected 0", len(allTasks))
	}
}
