package main

import (
	"bytes"
	"log"
	"os"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	os.Stdout.WriteString("Generating READMEs for quickstarts\n")

	const path = "/workspaces/terraform-provider-power-platform/examples/quickstarts/"
	quickstartTemplate, qerr := os.ReadFile("/workspaces/terraform-provider-power-platform/tools/quickstartgen/quickstart.md.tmpl")
	if qerr != nil {
		log.Fatal(qerr)
	}

	qsDir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer qsDir.Close()

	// Read directory entries
	fileInfos, err := qsDir.Readdir(-1) // Passing -1 to read all available entries
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	// Iterate over each entry and print subdirectories
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {

			modulePath := filepath.Join(path, fileInfo.Name())

			os.Stdout.WriteString("Generating README.md for " + modulePath + "\n")

			// Parse the terraform in the module
			module, diags := tfconfig.LoadModule(modulePath)
			if diags.HasErrors() {
				panic(diags)
			}

			// Parse the README.md.tmpl
			var readmeTemplate bytes.Buffer
			readmeTemplate.WriteString("<!-- This document is auto-generated. Do not edit directly. Make changes to README.md.tmpl instead. -->\n")
			err := RenderReadme(&readmeTemplate, filepath.Join(path, fileInfo.Name(), "README.md.tmpl"), struct {
				ModuleDetails string
				ModuleName    string
			}{
				ModuleDetails: string(quickstartTemplate),
				ModuleName:    fileInfo.Name(),
			})
			if err != nil {
				panic(err)
			}

			var readmeMarkdown bytes.Buffer
			RenderMarkdown(&readmeMarkdown, readmeTemplate.String(), module)

			os.WriteFile(filepath.Join(path, fileInfo.Name(), "README.md"), readmeMarkdown.Bytes(), 0644)
		}
	}
}

func RenderReadme(w io.Writer, templatePath string, data any) error {

	tmpl := template.Must(template.ParseFiles(templatePath))
	return tmpl.Execute(w, data)
}

func RenderMarkdown(w io.Writer, markdownTemplate string, module *tfconfig.Module) error {
	tmpl := template.New("md")
	tmpl.Funcs(template.FuncMap{
		"tt": func(s string) string {
			return "`" + s + "`"
		},
		"commas": func(s []string) string {
			return strings.Join(s, ", ")
		},
		"json": func(v interface{}) (string, error) {
			j, err := json.Marshal(v)
			return string(j), err
		},
		"severity": func(s tfconfig.DiagSeverity) string {
			switch s {
			case tfconfig.DiagError:
				return "Error: "
			case tfconfig.DiagWarning:
				return "Warning: "
			default:
				return ""
			}
		},
		"relativePath": func(path string) string {
			rp, _ := filepath.Rel("/workspaces/terraform-provider-power-platform", path)
			return rp
		},
	})
	template.Must(tmpl.Parse(markdownTemplate))
	return tmpl.Execute(w, module)
}

const detailsTemplate = `
<!-- This document is auto-generated. Do not edit directly. -->
# Module {{ .Path | relativePath | tt}}

{{- if .RequiredCore}}

## Core Version Constraints:

{{- range .RequiredCore }}
* {{ tt . }}
{{- end}}{{end}}

{{- if .RequiredProviders}}

## Provider Requirements:
{{- range $name, $req := .RequiredProviders }}
* **{{ $name }}{{ if $req.Source }} ({{ $req.Source | tt }}){{ end }}:** {{ if $req.VersionConstraints }}{{ commas $req.VersionConstraints | tt }}{{ else }}(any version){{ end }}
{{- end}}{{end}}

{{- if .Variables}}

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
{{- range .Variables }}
| {{ tt .Name }} | {{ .Description }} | | {{ json .Default | tt }} | {{ .Required }} |
{{- end}}
{{end}}

{{- if .Outputs}}

## Output Values

| Name | Description |
|------|-------------|
{{- range .Outputs }}
| {{ tt .Name }} | {{ .Description }} |
{{- end}}

{{end}}

{{- if .ManagedResources}}

## Managed Resources
{{- range .ManagedResources }}
* {{ printf "%s.%s" .Type .Name | tt }} from {{ tt .Provider.Name }}
{{- end}}{{end}}

{{- if .DataResources}}

## Data Resources
{{- range .DataResources }}
* {{ printf "data.%s.%s" .Type .Name | tt }} from {{ tt .Provider.Name }}
{{- end}}{{end}}

{{- if .ModuleCalls}}

## Child Modules
{{- range .ModuleCalls }}
* {{ tt .Name }} from {{ tt .Source }}{{ if .Version }} ({{ tt .Version }}){{ end }}
{{- end}}{{end}}

{{- if .Diagnostics}}

## Problems
{{- range .Diagnostics }}

## {{ severity .Severity }}{{ .Summary }}{{ if .Pos }}

(at {{ tt .Pos.Filename }} line {{ .Pos.Line }}{{ end }})
{{ if .Detail }}
{{ .Detail }}
{{- end }}

{{- end}}{{end}}

`
