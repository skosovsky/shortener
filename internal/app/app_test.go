package app_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/config"
	"shortener/internal/app"
	"shortener/internal/service"
	"shortener/internal/store"
)

func TestRoutingGet(t *testing.T) {
	t.Parallel()

	type want struct {
		path   string
		status int
		body   string
		// answer string // TODO: add later
	}

	testCases := []struct {
		name string
		want want
	}{
		{
			name: "with id",
			want: want{
				path:   "/10",
				status: 307,
				body:   "",
				// answer: "https://ya.ru", // TODO: add later
			},
		},
	}

	// experiments - TODO: use mock
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	cfg, _ := config.NewConfig()
	db, _ := store.NewDummyStore()
	shortener := service.NewSiteService(db, cfg)
	ctx := context.WithValue(context.Background(), app.KeyServiceCtx{}, shortener)

	server := httptest.NewUnstartedServer(app.Handler())
	server.Config.ConnContext = func(_ context.Context, _ net.Conn) context.Context { return ctx }
	server.Start()

	t.Cleanup(server.Close)

	// experiments - TODO: put to table test
	// db.EXPECT().Get("10").Return(model.Site{
	//	ID:        "10",
	//	Link:      "https://ya.ru",
	//	ShortLink: "http://localhost:8080/Jfdf00",
	// }, true)

	for _, tt := range testCases {
		tt := tt //nolint:copyloopvar // it's for stupid Yandex Practicum static test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			url := server.URL + tt.want.path
			request, err := http.NewRequest(http.MethodGet, url, http.NoBody)
			require.NoError(t, err)

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, response.StatusCode)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			err = response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(responseBody))
		})
	}
}
