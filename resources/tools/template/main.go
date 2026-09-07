package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateData struct {
	NAMESPACE              string
	MIGRATION_VERSION      string
	MIGRATION_NAME_VERSION string
}

func main() {
	src := flag.String("src", "", "Source template file")
	dst := flag.String("dst", "", "Destination file")

	namespace := flag.String("namespace", "local", "Kubernetes namespace")
	version := flag.String("version", "", "Migration version")

	flag.Parse()

	if *src == "" {
		fatal("missing required argument: -src")
	}

	if *dst == "" {
		fatal("missing required argument: -dst")
	}

	if *version == "" {
		fatal("missing required argument: -version")
	}

	// fmt.Printf("src = %q\n", *src)
	// fmt.Printf("dst = %q\n", *dst)

	data := TemplateData{
		NAMESPACE:              *namespace,
		MIGRATION_VERSION:      *version,
		MIGRATION_NAME_VERSION: sanitizeKubernetesName(*version),
	}

	if err := render(*src, *dst, data); err != nil {
		fatal("%v", err)
	}

	fmt.Printf("Rendered %s -> %s\n", *src, *dst)
}

func render(src, dst string, data TemplateData) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read template %q: %w", src, err)
	}

	tmpl, err := template.New(filepath.Base(src)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse template %q: %w", src, err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("create destination directory: %w", err)
	}

	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination %q: %w", dst, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	return nil
}

func sanitizeKubernetesName(version string) string {
	version = strings.ToLower(version)

	var builder strings.Builder

	for _, r := range version {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)

		case r >= '0' && r <= '9':
			builder.WriteRune(r)

		case r == '_':
			builder.WriteRune(r)

		default:
			builder.WriteRune('_')
		}
	}

	result := strings.Trim(builder.String(), "_")

	return result
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "template: "+format+"\n", args...)
	os.Exit(1)
}
