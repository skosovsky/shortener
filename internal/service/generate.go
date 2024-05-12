package service

import (
	"crypto/rand"

	"shortener/config"
	"shortener/internal/model"
)

const (
	LettersNums = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	LenID       = 8
)

func SiteGenerate(link string) model.Site {
	var idxDomain = 0
	var id = make([]byte, LenID)
	var site model.Site

	_, _ = rand.Read(id)

	for k, v := range id {
		id[k] = LettersNums[v%byte(len(LettersNums))]
	}

	domains, _ := config.GetDomains()
	domain := domains[idxDomain]

	site.ID = string(id)
	site.Link = link
	site.ShortLink = domain + "/" + string(id)

	return site
}
