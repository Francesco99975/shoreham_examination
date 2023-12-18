package models

import "time"

type SEO struct {
	Description string
	Keywords    string
}
type Site struct {
	AppName  string
	Title    string
	Metatags SEO
	Year     int
}

func GetDefaultSite(title string) Site {
	return Site{
		AppName:  "Shoreham Examination",
		Title:    title,
		Metatags: SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}
}


