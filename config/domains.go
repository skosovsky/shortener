package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

func GetDomains() ([]string, error) { // TODO: Перенести в контекст, чтобы вызывалось 1 раз или в конфиг или удалить
	pathConfig := "config.json"
	fileConfig, err := os.ReadFile(pathConfig)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var domains []Domain
	err = json.Unmarshal(fileConfig, &domains)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling domains: %w", err)
	}

	if len(domains) == 0 {
		domains = append(domains, Domain{URL: "http://localhost:8080"})
	}

	var domainNames []string
	for i := range domains {
		if err = domains[i].validate(); err == nil {
			domainNames = append(domainNames, domains[i].URL)
		}
	}

	return domainNames, nil
}

func (d Domain) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(d)
	if err != nil {
		return fmt.Errorf("failed to validate config %v: %w", d, err)
	}

	return nil
}
