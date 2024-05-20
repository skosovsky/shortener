package service

import (
	"crypto/rand"

	"shortener/internal/model"
)

type (
	IDGenerator struct {
		symbols string
		length  int
	}

	FakeIDGenerator struct {
		symbols string
		length  int
	}
)

func NewIDGenerator() IDGenerator {
	const idLength = 8

	return IDGenerator{
		symbols: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		length:  idLength,
	}
}

func (i IDGenerator) Generate(domain string, link string) model.Site {
	var id = make([]byte, i.length)
	var site model.Site

	_, _ = rand.Read(id)

	for k, v := range id {
		id[k] = i.symbols[v%byte(len(i.symbols))]
	}

	site.ID = string(id)
	site.Link = link
	site.ShortLink = domain + "/" + string(id)

	return site
}

func NewFakeIDGenerator() FakeIDGenerator {
	const idLength = 8

	return FakeIDGenerator{
		symbols: "aaaaaaaaaa",
		length:  idLength,
	}
}

func (i FakeIDGenerator) Generate(domain string, link string) model.Site {
	var id = make([]byte, i.length)
	var site model.Site

	_, _ = rand.Read(id)

	for k, v := range id {
		id[k] = i.symbols[v%byte(len(i.symbols))]
	}

	site.ID = string(id)
	site.Link = link
	site.ShortLink = domain + "/" + string(id)

	return site
}
