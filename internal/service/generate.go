package service

import (
	"crypto/rand"

	"shortener/internal/model"
)

const (
	LettersNums = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	LenID       = 8
)

func (s Shortener) SiteGenerate(link string) model.Site {
	var id = make([]byte, LenID)
	var site model.Site

	_, _ = rand.Read(id)

	for k, v := range id {
		id[k] = LettersNums[v%byte(len(LettersNums))]
	}

	site.ID = string(id)
	site.Link = link
	site.ShortLink = s.config.Domain.URL + "/" + string(id)

	return site
}
