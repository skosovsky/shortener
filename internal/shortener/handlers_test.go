package shortener_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/config"
	"shortener/internal/service"
	"shortener/internal/shortener"
	"shortener/internal/store"
)

func TestMethods(t *testing.T) {
	t.Parallel()

	type want struct {
		code        int
		response    string
		contentType string
	}
	testCases := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:    "Method Post",
			method:  http.MethodPost,
			request: "/",
			want: want{
				code:        201,
				response:    "/aaaaaaaa",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Method Get",
			method:  http.MethodGet,
			request: "/1",
			want: want{
				code:        307,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Method Put",
			method:  http.MethodPut,
			request: "/1",
			want: want{
				code:        405,
				response:    "Method Not Allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Method Delete",
			method:  http.MethodDelete,
			request: "/1",
			want: want{
				code:        405,
				response:    "Method Not Allowed\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	var cfg config.Config
	db := store.NewDummyStore()
	generator := service.NewFakeIDGenerator()
	shortenerService := service.NewService(db, cfg, generator)
	handler := shortener.NewHandler(shortenerService)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(tt.method, tt.request, http.NoBody)
			responseRecorder := httptest.NewRecorder()

			switch tt.method {
			case http.MethodPost:
				handler.AddSite(responseRecorder, request)
			case http.MethodGet:
				handler.GetSite(responseRecorder, request)
			default:
				handler.AddSite(responseRecorder, request)
			}

			response := responseRecorder.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(responseBody))
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
		})
	}
}
