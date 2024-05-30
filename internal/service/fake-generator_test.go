package service_test

import (
	"reflect"
	"testing"

	"shortener/internal/service"
)

func TestFakeIDGenerator_Generate(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			domain string
			link   string
		}

		request struct {
			args args
			want service.Site
		}
	)

	first := request{
		args: args{
			domain: "localhost:8080",
			link:   "test.ru",
		},
		want: service.Site{
			ID:        "byficVeg",
			Link:      "test.ru",
			ShortLink: "localhost:8080/byficVeg",
		},
	}

	second := request{
		args: args{
			domain: "localhost:8080",
			link:   "test.ru",
		},
		want: service.Site{
			ID:        "ufehHgNr",
			Link:      "test.ru",
			ShortLink: "localhost:8080/ufehHgNr",
		},
	}

	third := request{
		args: args{
			domain: "localhost:8080",
			link:   "test.ru",
		},
		want: service.Site{
			ID:        "VggrObRX",
			Link:      "test.ru",
			ShortLink: "localhost:8080/VggrObRX",
		},
	}

	generator := service.NewFakeIDGenerator()

	t.Run("generate 3 times", func(t *testing.T) {
		t.Parallel()

		if got := generator.Generate(first.args.domain, first.args.link); !reflect.DeepEqual(got, first.want) {
			t.Errorf("Generate() = %v, want %v", got, first.want)
		}

		if got := generator.Generate(second.args.domain, second.args.link); !reflect.DeepEqual(got, second.want) {
			t.Errorf("Generate() = %v, want %v", got, second.want)
		}

		if got := generator.Generate(third.args.domain, third.args.link); !reflect.DeepEqual(got, third.want) {
			t.Errorf("Generate() = %v, want %v", got, third.want)
		}
	})
}
