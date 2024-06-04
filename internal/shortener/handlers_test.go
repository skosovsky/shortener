package shortener_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/config"
	"shortener/internal/log"
	"shortener/internal/service"
	"shortener/internal/shortener"
	"shortener/internal/store"
)

func prepare(t *testing.T) {
	t.Helper()

	log.Prepare()
}

func TestPostAddSite(t *testing.T) { //TODO: Add new tests with compress
	prepare(t)

	t.Parallel()

	type want struct {
		code        int
		response    string
		contentType string
	}

	testCases := []struct {
		name        string
		requestBody string
		want        want
	}{
		{
			name:        "Add empty site",
			requestBody: "",
			want: want{
				code:        400,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Add not valid site",
			requestBody: "ya.ru",
			want: want{
				code:        400,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Add valid new site",
			requestBody: "http://ya.ru",
			want: want{
				code:        201,
				response:    "/byficVeg",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	var cfg config.Config
	db := store.NewDummyStore()
	generator := service.NewFakeIDGenerator()
	shortenerService := service.NewService(db, cfg, &generator)
	handler := shortener.NewHandler(shortenerService)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			responseRecorder := httptest.NewRecorder()

			handler.AddSite(responseRecorder, request)

			response := responseRecorder.Result()

			assert.Equal(t, tt.want.code, response.StatusCode)
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			err = response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(responseBody))
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
		})
	}
}

func TestPostAddDuplSite(t *testing.T) {
	prepare(t)

	t.Parallel()

	type (
		want struct {
			code     int
			response string
		}

		request struct {
			requestBody string
			want        want
		}
	)

	first := request{
		requestBody: "http://ya.ru",
		want: want{
			code:     201,
			response: "/byficVeg",
		},
	}

	second := request{
		requestBody: "http://ya.ru",
		want: want{
			code:     201,
			response: "/ufehHgNr",
		},
	}

	var cfg config.Config
	db := store.NewDummyStore()
	generator := service.NewFakeIDGenerator()
	shortenerService := service.NewService(db, cfg, &generator)
	handler := shortener.NewHandler(shortenerService)

	t.Run("Add two same site", func(t *testing.T) {
		t.Parallel()

		firstRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(first.requestBody))
		secondRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(second.requestBody))
		firstResponseRecorder := httptest.NewRecorder()
		secondResponseRecorder := httptest.NewRecorder()

		handler.AddSite(firstResponseRecorder, firstRequest)
		firstResponse := firstResponseRecorder.Result()

		handler.AddSite(secondResponseRecorder, secondRequest)
		secondResponse := secondResponseRecorder.Result()

		assert.Equal(t, first.want.code, firstResponse.StatusCode)
		assert.Equal(t, second.want.code, firstResponse.StatusCode)

		firstResponseBody, err := io.ReadAll(firstResponse.Body)
		require.NoError(t, err)
		secondResponseBody, err := io.ReadAll(secondResponse.Body)
		require.NoError(t, err)

		err = firstResponse.Body.Close()
		require.NoError(t, err)
		err = secondResponse.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, first.want.response, string(firstResponseBody))
		assert.Equal(t, second.want.response, string(secondResponseBody))
	})
}

func TestPostAddSiteJSON(t *testing.T) {
	prepare(t)

	t.Parallel()

	type want struct {
		code        int
		response    string
		contentType string
	}

	testCases := []struct {
		name        string
		requestBody string
		want        want
	}{
		{
			name:        "Add empty body",
			requestBody: "",
			want: want{
				code:        400,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Add empty JSON",
			requestBody: `{}`,
			want: want{
				code:        400,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Add empty JSON 2",
			requestBody: `{"":""}`,
			want: want{
				code:        400,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Add not valid JSON",
			requestBody: `{"url":"ya.ru"`,
			want: want{
				code:        400,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Add not valid site, but valid JSON",
			requestBody: `{"url":"ya.ru"}`,
			want: want{
				code:        400,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Add valid new site and valid JSON",
			requestBody: `{"url":"http://ya.ru"}`,
			want: want{
				code:        201,
				response:    `{"result":"/byficVeg"}` + "\n",
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	var cfg config.Config
	db := store.NewDummyStore()
	generator := service.NewFakeIDGenerator()
	shortenerService := service.NewService(db, cfg, &generator)
	handler := shortener.NewHandler(shortenerService)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			responseRecorder := httptest.NewRecorder()

			handler.AddSiteJSON(responseRecorder, request)

			response := responseRecorder.Result()

			assert.Equal(t, tt.want.code, response.StatusCode)
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			err = response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(responseBody))
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
		})
	}
}

func TestIsValidURL(t *testing.T) {
	prepare(t)

	t.Parallel()

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "empty url",
			url:  "",
			want: false,
		},
		{
			name: "wrong url",
			url:  "ya",
			want: false,
		},
		{
			name: "wrong url 2",
			url:  "ya.ru",
			want: false,
		},
		{
			name: "valid url",
			url:  "http://ya.ru",
			want: true,
		},
		{
			name: "wrong url 2",
			url:  "https://ya.ru",
			want: true,
		},
		{
			name: "wrong url 3",
			url:  "http://ya.ru:8080",
			want: true,
		},
	}

	var cfg config.Config
	db := store.NewDummyStore()
	generator := service.NewFakeIDGenerator()
	shortenerService := service.NewService(db, cfg, &generator)
	handler := shortener.NewHandler(shortenerService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, handler.IsValidURL(tt.url), "isValidURL(%v)", tt.url)
		})
	}
}
