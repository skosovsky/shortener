package shortener

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"

	"shortener/internal/log"
	"shortener/internal/mux"
	"shortener/internal/service"
)

type Handler struct {
	service service.Shortener
}

func NewHandler(service service.Shortener) Handler {
	return Handler{service: service}
}

func (h Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.Use(WithLogging)
	router.Use(WithGzipCompress)

	router.Post("/", h.AddSite)
	router.Post("/api/shorten", h.AddSiteJSON)
	router.Get("/{id}", h.GetSite)
	router.Get("/ping", h.Ping)

	return router
}

func (h Handler) AddSiteJSON(w http.ResponseWriter, r *http.Request) { //TODO: Update swagger api docs
	url := struct {
		URL string `json:"url"`
	}{
		URL: "",
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("error writing response", //nolint:contextcheck // no ctx
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	if len(body) == 0 {
		log.Debug("empty body") //nolint:contextcheck // no ctx

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	err = json.Unmarshal(body, &url)
	if err != nil {
		log.Debug("error decode to json", //nolint:contextcheck // no ctx
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	defer func(Body io.ReadCloser) { //nolint:contextcheck // no ctx
		err = Body.Close()
		if err != nil {
			log.Error("error close body",
				log.ErrAttr(err))
		}
	}(r.Body)

	if !h.IsValidURL(url.URL) {
		log.Debug("url validate failed") //nolint:contextcheck // no ctx

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	site, err := h.service.Add(url.URL)
	if err != nil {
		log.Error("error site add", //nolint:contextcheck // no ctx
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	result := struct {
		URL string `json:"result"`
	}{
		URL: site.ShortLink,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error("error encode to json", //nolint:contextcheck // no ctx
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h Handler) AddSite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("error reading response", //nolint:contextcheck // no ctx
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	defer func(Body io.ReadCloser) { //nolint:contextcheck // no ctx
		err = Body.Close()
		if err != nil {
			log.Error("error close body",
				log.ErrAttr(err))
		}
	}(r.Body)

	if !h.IsValidURL(string(body)) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	site, err := h.service.Add(string(body))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_, err = io.WriteString(w, site.ShortLink)
	if err != nil {
		log.Error("error writing response", //nolint:contextcheck // no ctx
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h Handler) GetSite(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	site, err := h.service.Get(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Location", site.Link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	err := h.service.Ping()
	if err != nil {
		log.Error("error pinging database", //nolint:contextcheck // no ctx
			log.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (Handler) IsValidURL(url string) bool {
	if len(url) == 0 {
		return false
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	checkable := struct {
		URL string `validate:"url"`
	}{
		URL: url,
	}

	err := validate.Struct(checkable)

	return err == nil
}
