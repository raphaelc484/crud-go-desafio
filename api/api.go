package api

import (
	"crud-go-desafio/db"
	"crud-go-desafio/model"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

var database = db.NewApplication()

func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/api/users", handlerInsert())
	r.Get("/api/users", handlerFindAll())
	r.Get("/api/users/{id}", handlerFindById())
	r.Put("/api/users/{id}", handlerUpdate())
	r.Delete("/api/users/{id}", handlerDelete())

	return r
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("error ao fazer marshal de json", "error", err)
		sendJSON(
			w,
			Response{Error: "something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("error ao enviar a resposta", "error", err)
	}
}

func handlerDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			sendJSON(
				w,
				Response{Error: "Method not allowed"},
				http.StatusMethodNotAllowed,
			)
			return
		}

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			sendJSON(
				w,
				Response{Error: "Invalid ID"},
				http.StatusBadRequest,
			)
			return
		}

		found := database.Delete(uuid)
		if !found {
			sendJSON(
				w,
				Response{Error: "The user with the specified ID does not exist"},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(
			w,
			Response{Data: "User Deleted"},
			http.StatusOK,
		)
	}
}

func handlerUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			sendJSON(
				w,
				Response{Error: "Method not allowed"},
				http.StatusMethodNotAllowed,
			)
			return
		}

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			sendJSON(
				w,
				Response{Error: "Invalid ID"},
				http.StatusBadRequest,
			)
			return
		}

		var req model.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(
				w,
				Response{Error: "Invalid request payload"},
				http.StatusBadRequest,
			)
			return
		}

		user := req.ToUser()
		if err := user.Validate(); err != nil {
			sendJSON(
				w,
				Response{Error: err.Error()},
				http.StatusBadRequest,
			)
			return
		}

		updatedUser, found := database.Update(uuid, req.FirstName, req.LastName, req.Biography)
		if !found {
			sendJSON(
				w,
				Response{Error: "The user with the specified ID does not exist"},
				http.StatusNotFound,
			)
			return
		}

		sendJSON(
			w,
			Response{Data: updatedUser},
			http.StatusOK,
		)
	}
}

func handlerFindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sendJSON(
				w,
				Response{Error: "Method not allowed"},
				http.StatusMethodNotAllowed,
			)
			return
		}

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			sendJSON(
				w,
				Response{Error: "Invalid Id"},
				http.StatusBadRequest,
			)
			return
		}

		user, found := database.FindById(uuid)
		if !found {
			sendJSON(
				w,
				Response{Error: "The user with the specified ID does not exist"},
				http.StatusBadRequest,
			)
			return
		}

		sendJSON(
			w,
			Response{Data: user},
			http.StatusOK,
		)
	}
}

func handlerFindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sendJSON(
				w,
				Response{Error: "Method not allowed"},
				http.StatusMethodNotAllowed,
			)
			return
		}
		users := database.FindAll()
		sendJSON(
			w,
			Response{Data: users},
			http.StatusOK,
		)
	}
}

func handlerInsert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendJSON(
				w,
				Response{Error: "Method not allowed"},
				http.StatusMethodNotAllowed,
			)
			return
		}

		var req model.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(
				w,
				Response{Error: "Invalid request payload"},
				http.StatusMethodNotAllowed,
			)
			return
		}

		user := req.ToUser()
		if err := user.Validate(); err != nil {
			sendJSON(
				w,
				Response{Error: err.Error()},
				http.StatusBadRequest,
			)
			return
		}

		newUser := database.Insert(req.FirstName, req.LastName, req.Biography)
		sendJSON(
			w,
			Response{Data: newUser},
			http.StatusCreated,
		)
	}
}
