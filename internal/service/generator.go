package service

import (
	"crypto/rand"
)

type IDGenerator struct {
	symbols string
	length  int
}

func NewIDGenerator() IDGenerator {
	const idLength = 8

	return IDGenerator{
		symbols: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		length:  idLength,
	}
}

func (i IDGenerator) Generate(domain string, link string) Site {
	var id = make([]byte, i.length)
	var site Site

	_, _ = rand.Read(id)

	for k, v := range id {
		id[k] = i.symbols[v%byte(len(i.symbols))]
	}

	site.ID = string(id)
	site.Link = link
	site.ShortLink = domain + "/" + string(id)

	return site
}
