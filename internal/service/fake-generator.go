package service

import (
	"math/rand"
)

type FakeIDGenerator struct {
	symbols string
	length  int
	counter int64
}

func NewFakeIDGenerator() FakeIDGenerator {
	const idLength = 8

	return FakeIDGenerator{
		symbols: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		length:  idLength,
		counter: 0,
	}
}

func (i *FakeIDGenerator) Generate(domain string, link string) Site {
	random := rand.New(rand.NewSource(i.counter)) //nolint:gosec // fake
	i.counter++

	var id = make([]byte, i.length)
	var site Site

	_, _ = random.Read(id)

	for k, v := range id {
		id[k] = i.symbols[v%byte(len(i.symbols))]
	}

	site.ID = string(id)
	site.Link = link
	site.ShortLink = domain + "/" + string(id)

	return site
}
