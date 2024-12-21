package models

type Route struct {
	Params  map[string]string
	Handler string
}
