package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/hyperboloide/pipe"
	"github.com/hyperboloide/pipe/rw"
	"github.com/segmentio/ksuid"
)

type ServiceDefinition struct {
	URL        string          `json:"url"`
	WriterPipe json.RawMessage `json:"writer,omitempty"`
	ReaderPipe json.RawMessage `json:"reader,omitempty"`
	Deleter    json.RawMessage `json:"deleter,omitempty"`
}

func (sd *ServiceDefinition) Error(err error) {
	log.Fatalf("service '%s' encountred an error: %s", sd.URL, err)
}

func SetHandler(r *chi.Mux, sd ServiceDefinition) {
	r.Route("/"+sd.URL, func(r chi.Router) {

		if sd.ReaderPipe != nil {
			if ops, err := NewReadOperationsFromJson(sd.ReaderPipe); err != nil {
				sd.Error(err)
			} else {
				SetReadHandler(r, ops)
			}
		}

		if sd.WriterPipe != nil {
			if ops, err := NewWriteOperationsFromJson(sd.WriterPipe); err != nil {
				sd.Error(err)
			} else {
				SetWriteHandler(r, ops)
			}
		}

		if sd.Deleter != nil {
			if del, err := DeleterFromJson(sd.Deleter); err != nil {
				sd.Error(err)
			} else {
				SetDeleteHandler(r, del)
			}
		}

	})

}

func SetReadHandler(r chi.Router, ops *ReadOperations) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if reader, err := ops.Input.NewReader(id); err != nil {
			http.Error(w, http.StatusText(404), 404)
		} else {
			defer reader.Close()
			p := pipe.New(reader)
			if err := ops.SetPipe(p); err != nil {
				http.Error(w, http.StatusText(500), 500)
			}
			p.To(w)
			if err := p.Exec(); err != nil {
				http.Error(w, http.StatusText(500), 500)
			}
		}
	}

	r.Get("/{id}", handler)
}

func SetDeleteHandler(r chi.Router, del rw.Deleter) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := del.Delete(id); err != nil {
			http.Error(w, http.StatusText(500), 500)
		} else {
			http.Error(w, http.StatusText(204), 204)
		}
	}

	r.Delete("/{id}", handler)
}

func SetWriteHandler(r chi.Router, ops *WriteOperations) {
	handler := func(id string, w http.ResponseWriter, r *http.Request) {
		var reader io.Reader
		if fr, _, err := r.FormFile("file"); err != nil {
			defer r.Body.Close()
			reader = r.Body
		} else {
			defer fr.Close()
			reader = fr
		}
		p := pipe.New(reader)
		if err := ops.SetPipe(p, id); err != nil {
			http.Error(w, http.StatusText(500), 500)
		} else if err := p.Exec(); err != nil {
			http.Error(w, http.StatusText(500), 500)
		} else {
			data := struct {
				In  int64 `bytes_in`
				Out int64 `bytes_out`
			}{p.TotalIn, p.TotalOut}
			if res, err := json.Marshal(data); err != nil {
				http.Error(w, http.StatusText(500), 500)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(201)
				if _, err = w.Write(res); err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
			}
		}
	}

	generateID := func(w http.ResponseWriter, r *http.Request) {
		id := ksuid.New().String()
		handler(id, w, r)
	}

	extractID := func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		handler(id, w, r)
	}

	r.Post("/", generateID)
	r.Put("/", generateID)
	r.Post("/{id}", extractID)
	r.Put("/{id}", extractID)
}
