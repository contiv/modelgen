package generators

import (
	"bytes"
	"text/template"

	"github.com/contiv/modelgen/texthelpers"
)

var templateMap = map[string]*template.Template{}

var funcMap = template.FuncMap{
	"initialCap":        texthelpers.InitialCap,
	"initialLow":        texthelpers.InitialLow,
	"depunct":           texthelpers.Depunct,
	"capFirst":          texthelpers.CapFirst,
	"translateCfgType":  texthelpers.TranslateCfgPropertyType,
	"translateOperType": texthelpers.TranslateOperPropertyType,
}

func ParseTemplates() error {
	for name, content := range templates {
		var err error
		templateMap[name], err = template.New(name).Funcs(funcMap).Parse(content)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetTemplate(templateName string) *template.Template {
	return templateMap[templateName]
}

func RunTemplate(templateName string, obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl := GetTemplate(templateName)
	if err := tmpl.Execute(buf, obj); err != nil {
		return nil, err
	}

	return bytes.TrimSuffix(buf.Bytes(), []byte("\n  ")), nil
}
