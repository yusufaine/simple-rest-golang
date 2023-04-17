package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"yusufaine/golang-todo/internal/taskstore"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	return &taskServer{
		store: taskstore.New(),
	}
}

func (ts *taskServer) taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/task" || r.URL.Path == "/task/" {
		switch r.Method {
		case http.MethodGet:
			ts.getAllTasksHandler(w, r)
		case http.MethodPost:
			ts.createTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	} else {

		// /task/<id> handler
		path := strings.Trim(r.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) != 2 {
			http.Error(w, "Bad request, expected /task/<id>", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			ts.getTaskHandler(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, expected GET /tag/<tag>", http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 2 {
		http.Error(w, "Bad request, expected /tag/<tag>", http.StatusBadRequest)
		return
	}

	tag := pathParts[1]
	tasks := ts.store.GetTasksByTag(tag)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})
	res, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, expected GET /due/<date>", http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	badRequestError := func() {
		http.Error(w, "Bad request, expected /due/<year>/<month>/<day>", http.StatusBadRequest)
	}

	if len(pathParts) != 4 {
		badRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])
	if err != nil {
		badRequestError()
		return
	}

	month, err := strconv.Atoi(pathParts[2])
	if err != nil || month < (int)(time.January) || month > (int)(time.December) {
		badRequestError()
		return
	}

	day, err := strconv.Atoi(pathParts[3])
	if err != nil {
		badRequestError()
		return
	}

	tasks := ts.store.GetTasksByDate(year, time.Month(month), day)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})
	res, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	allTasks := ts.store.GetAllTasks()
	sort.Slice(allTasks, func(i, j int) bool {
		return allTasks[i].Id < allTasks[j].Id
	})
	res, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	type createTaskRequest struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type createTaskResponse struct {
		Id int `json:"id"`
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	var req createTaskRequest
	if err := decode.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(req.Text, req.Tags, req.Due)
	res, err := json.Marshal(createTaskResponse{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	mux := http.NewServeMux()
	ts := NewTaskServer()
	mux.HandleFunc("/task", ts.taskHandler)
	mux.HandleFunc("/task/", ts.taskHandler)
	mux.HandleFunc("/tag", ts.tagHandler)
	mux.HandleFunc("/tag/", ts.tagHandler)
	mux.HandleFunc("/due", ts.dueHandler)
	mux.HandleFunc("/due/", ts.dueHandler)

	// Start server and log running port
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
