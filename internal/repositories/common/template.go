package common

import (
	"bytes"
	"fmt"
	"html/template"
)

// ProcessTemplate is responsible to fill a template text with data from the provided Struct.
func ProcessTemplate(txt string, obj any) (string, error) {
	t, err := template.New("Template").Parse(txt)
	if err != nil {
		return "", fmt.Errorf("ProcessTemplate Parse failed: %w", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, obj); err != nil {
		return "", fmt.Errorf("ProcessTemplate Execute failed: %w", err)
	}

	return tpl.String(), nil
}
