package models

import "github.com/elliotchance/orderedmap/v3"

type Route struct {
	Params  *orderedmap.OrderedMap[string, string]
	Handler string
}
