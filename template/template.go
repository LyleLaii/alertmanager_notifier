package template

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/pkg/utils"
	"bytes"
	"fmt"
	tmpltext "text/template"
)

type Template struct {
	tmpltext.Template
}

func GenTempInstance(tmplFilePath string) (*Template, error) {
	var tmpl *tmpltext.Template
	var err error
	funcMap := tmpltext.FuncMap{"uuid": utils.NewUUID,
		"uuid32": utils.NewUUID32,
		"transTime": utils.TransTimeZoneAutoCustom}
	if utils.FileExist(tmplFilePath) {
		tmpl, err = tmpltext.ParseFiles(tmplFilePath)
		if err != nil {
			return &Template{}, err
		}
	} else {
		tmpl = tmpltext.New("")
	}
	tmpl = tmpl.Funcs(funcMap).Option("missingkey=zero")

	return &Template{*tmpl}, nil
}

func generateAlertTmplVariable(index int) string {
	tmplString := fmt.Sprintf(config.AlertsTmplVariable, index)
	for _, value := range config.AlertInfoTmplVariable {
		tmplString += value
	}

	return tmplString
}

func (t *Template) ParseTmplString(alertIndex int, tmplString string, data interface{}) (string, error) {
	tmpl, err := t.Clone()
	if err != nil {
		return "", err
	}
	alertTmplString := generateAlertTmplVariable(alertIndex)
	tmplStr := alertTmplString + tmplString
	tmpl, err = tmpl.Parse(tmplStr)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}