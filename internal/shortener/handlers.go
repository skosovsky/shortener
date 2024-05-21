package shortener

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"shortener/internal/log"
	"shortener/internal/service"
)

type Handler struct {
	service service.Shortener
}

func NewHandler(service service.Shortener) Handler {
	return Handler{service: service}
}

func (h Handler) InitRoutes() http.Handler {
	router := chi.NewRouter()

	router.Post("/", h.AddSite)
	router.Get("/{id}", h.GetSite)

	return router
}

func (h Handler) AddSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Error writing response", //nolint:contextcheck // false positive
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()

	site, err := h.service.Add(string(body))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_, err = io.WriteString(w, site.ShortLink)
	if err != nil {
		log.Error("Error writing response", //nolint:contextcheck // false positive
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h Handler) GetSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	id := chi.URLParam(r, "id")

	site, err := h.service.Get(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Location", site.Link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
