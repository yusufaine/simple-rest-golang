package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextId: 0,
	}
}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks[ts.nextId] = Task{
		Id:   ts.nextId,
		Text: text,
		Tags: tags,
		Due:  due,
	}

	ts.nextId++

	return ts.nextId - 1
}

func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, ok := ts.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task with id %d not found", id)
	}

	return task, nil
}

func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				break
			}
		}
	}
	return tasks
}

func (ts *TaskStore) GetTasksByDate(year int, month time.Month, day int) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	tasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}
