{{ define "test" }}Test{{ end }}
{{ define "test1" }}{{if eq . "firing"}}T1{{else}}T0{{end}}{{ end }}
{{ define "severity" }}{{if eq . "error"}}E{{else if eq . "warn"}}W{{else}}I{{end}}{{ end }}