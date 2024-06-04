package shortener_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/config"
	"shortener/internal/service"
	"shortener/internal/shortener"
	"shortener/internal/store"
)

func TestRouting(t *testing.T) {
	prepare(t)

	t.Parallel()

	type want struct {
		status int
	}

	testCases := []struct {
		name        string
		method      string
		request     string
		requestBody string
		want        want
	}{
		{
			name:        "Post with request, empty body",
			method:      http.MethodPost,
			request:     "/ping",
			requestBody: "",
			want: want{
				status: 400,
			},
		},
		{
			name:        "Post with request, valid body",
			method:      http.MethodPost,
			request:     "/ping",
			requestBody: "http://ya.ru",
			want: want{
				status: 201,
			},
		},
		{
			name:        "Post without ID, valid body",
			method:      http.MethodPost,
			request:     "/",
			requestBody: "http://ya.ru",
			want: want{
				status: 201,
			},
		},
	}

	var cfg config.Config
	db := store.NewDummyStore()
	generator := service.NewFakeIDGenerator()
	shortenerService := service.NewService(db, cfg, &generator)
	handler := shortener.NewHandler(shortenerService)

	server := httptest.NewServer(handler.InitRoutes())

	t.Cleanup(server.Close)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			url := server.URL + tt.request
			request, err := http.NewRequest(tt.method, url, strings.NewReader(tt.requestBody))
			require.NoError(t, err)

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, response.StatusCode)

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}
