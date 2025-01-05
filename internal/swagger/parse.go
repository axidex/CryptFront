package swagger

import (
	"encoding/json"
	"front/internal/models"
	"strconv"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/go-openapi/loads"
)

func GetRoutes(data json.RawMessage) (map[string]models.Route, error) {
	routes := make(map[string]models.Route)
	openApi, err := loads.Analyzed(data, "2.0")
	if err != nil {
		return nil, err
	}

	for routePath, path := range openApi.Spec().Paths.Paths {
		var route models.Route
		params := orderedmap.NewOrderedMap[string, string]()
		route.Handler = routePath

		for _, parameter := range path.Post.Parameters {
			var param string
			param, ok := parameter.Default.(string)
			if !ok {
				param = strconv.Itoa(int(parameter.Default.(float64)))
			}
			params.Set(parameter.Name, param)
		}
		route.Params = params

		routes[path.Post.Description] = route
	}

	return routes, nil
}
