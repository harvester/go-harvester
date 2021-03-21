package generator

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	types "github.com/rancher/wrangler/pkg/schemas"
)

var (
	blackListTypes = map[string]bool{
		"schema":     true,
		"resource":   true,
		"collection": true,
	}
)

func generateType(outputDir string, schema *types.Schema, schemas *types.Schemas) error {
	filePath := strings.ToLower("zz_generated_" + addUnderscore(schema.ID) + ".go")
	output, err := os.Create(path.Join(outputDir, filePath))
	if err != nil {
		return err
	}
	defer output.Close()

	tpl, err := template.New("type.template").
		Funcs(funcs()).
		Parse(strings.Replace(typeTemplate, "%BACK%", "`", -1))
	if err != nil {
		return err
	}

	return tpl.Execute(output, map[string]interface{}{
		"schema":            schema,
		"resourceActions":   getResourceActions(schema, schemas),
		"collectionActions": getCollectionActions(schema, schemas),
	})
}

func generateClient(outputDir string, schemas []*types.Schema) error {
	tpl, err := template.New("client.template").
		Funcs(funcs()).
		Parse(clientTemplate)
	if err != nil {
		return err
	}

	output, err := os.Create(path.Join(outputDir, "zz_generated_client.go"))
	if err != nil {
		return err
	}
	defer output.Close()

	return tpl.Execute(output, map[string]interface{}{
		"schemas": schemas,
	})
}

func GenerateClient(schemas *types.Schemas, privateTypes map[string]bool, outputDir, cattleOutputPackage string) error {
	baseDir := DefaultSourceTree()
	cattleDir := path.Join(outputDir, cattleOutputPackage)

	if err := prepareDirs(cattleDir); err != nil {
		return err
	}

	var cattleClientTypes []*types.Schema
	for _, schema := range schemas.Schemas() {
		if blackListTypes[schema.ID] {
			continue
		}

		if err := generateType(cattleDir, schema, schemas); err != nil {
			return err
		}

		if _, privateType := privateTypes[schema.ID]; !privateType {
			cattleClientTypes = append(cattleClientTypes, schema)
		}
	}

	if err := generateClient(cattleDir, cattleClientTypes); err != nil {
		return err
	}

	return Gofmt(baseDir, filepath.Join(outputDir, cattleOutputPackage))
}
