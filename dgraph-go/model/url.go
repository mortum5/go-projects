package model

import "time"

type Url struct {
	Uid    string    `json:"uid,omitempty"`
	URL    string    `json:"url,omitempty"`
	Slug   string    `json:"slug,omitempty"`
	Expire time.Time `json:"expire,omitempty"`
}

func New(url, slug string, d time.Duration) Url {
	return Url{
		URL:    url,
		Slug:   slug,
		Expire: time.Now().Add(d),
	}
}
