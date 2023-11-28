package models

type SEO struct {
	Description string
	Keywords string
}
type Site struct {
	AppName string
	Title string
	Metatags SEO
	CSRF string
	Year int
}