package generator

import (
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/rancher/wrangler/pkg/data/convert"
	types "github.com/rancher/wrangler/pkg/schemas"
	"github.com/rancher/wrangler/pkg/slice"
)

var (
	underscoreRegexp = regexp.MustCompile(`([a-z])([A-Z])`)
)

func funcs() template.FuncMap {
	return template.FuncMap{
		"capitalize":          convert.Capitalize,
		"unCapitalize":        convert.Uncapitalize,
		"upper":               strings.ToUpper,
		"toLower":             strings.ToLower,
		"hasGet":              hasGet,
		"hasPost":             hasPost,
		"getCollectionOutput": getCollectionOutput,
		"namespaced":          namespaced,
	}
}

func addUnderscore(input string) string {
	return strings.ToLower(underscoreRegexp.ReplaceAllString(input, `${1}_${2}`))
}

func hasGet(schema *types.Schema) bool {
	return slice.ContainsString(schema.CollectionMethods, http.MethodGet)
}

func namespaced(schema *types.Schema) bool {
	return schema.Attributes["namespaced"].(bool)
}

func hasPost(schema *types.Schema) bool {
	return slice.ContainsString(schema.CollectionMethods, http.MethodPost)
}

func getCollectionOutput(output, codeName string) string {
	if output == "collection" {
		return codeName + "Collection"
	}
	return convert.Capitalize(output)
}

func getResourceActions(schema *types.Schema, schemas *types.Schemas) map[string]types.Action {
	result := map[string]types.Action{}
	for name, action := range schema.ResourceActions {
		if action.Output != "" {
			if schemas.Schema(action.Output) != nil {
				result[name] = action
			}
		} else {
			result[name] = action
		}
	}
	return result
}

func getCollectionActions(schema *types.Schema, schemas *types.Schemas) map[string]types.Action {
	result := map[string]types.Action{}
	for name, action := range schema.CollectionActions {
		if action.Output != "" {
			output := action.Output
			if action.Output == "collection" {
				output = strings.ToLower(schema.CodeName)
			}
			if schemas.Schema(output) != nil {
				result[name] = action
			}
		} else {
			result[name] = action
		}
	}
	return result
}
