package model

type Site struct {
	ID        string
	Link      string
	ShortLink string
}

func (s Site) String() string {
	return s.Link + " > " + s.ShortLink
}
