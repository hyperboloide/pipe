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

// Definition defines a URL enpoint and at lease a either a Writer,
// a Reader or a Deleter.
type Definition struct {
	URL        string          `json:"url"`
	WriterPipe json.RawMessage `json:"writer,omitempty"`
	ReaderPipe json.RawMessage `json:"reader,omitempty"`
	Deleter    json.RawMessage `json:"deleter,omitempty"`
}

// Error display the service url and the encountred error.
func (d *Definition) Error(err error) {
	log.Fatalf("service '%s' encountred an error: %s", d.URL, err)
}

// SetHandler set the right handler in chi for the provoded ServiceDefinition.
func SetHandler(r *chi.Mux, d Definition) {
	r.Route("/"+d.URL, func(r chi.Router) {

		if d.ReaderPipe != nil {
			if ops, err := NewReadOperationsFromJSON(d.ReaderPipe); err != nil {
				d.Error(err)
			} else {
				SetReadHandler(r, ops)
			}
		}

		if d.WriterPipe != nil {
			if ops, err := NewWriteOperationsFromJSON(d.WriterPipe); err != nil {
				d.Error(err)
			} else {
				SetWriteHandler(r, ops)
			}
		}

		if d.Deleter != nil {
			if del, err := DeleterFromJSON(d.Deleter); err != nil {
				d.Error(err)
			} else {
				SetDeleteHandler(r, del)
			}
		}

	})

}

// SetReadHandler sets a chi handler for a Reader.
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

// SetDeleteHandler sets a chi handler for a Deleter.
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

// WriteResponse is returned as a json response on a sucessfull write.
type WriteResponse struct {
	ID       string `json:"id"`
	BytesIn  int64  `json:"bytes_in"`
	BytesOut int64  `json:"bytes_out"`
}

// SetWriteHandler sets a chi router for a Writer.
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
			data := &WriteResponse{id, p.TotalIn, p.TotalOut}

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
