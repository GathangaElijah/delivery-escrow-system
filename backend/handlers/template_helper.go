package handlers

import (
    "fmt"
    "html/template"
    "os"
    "path/filepath"
)

// GetTemplate loads a template with common functions
func GetTemplate(templateName string) (*template.Template, error) {
    // Define common template functions
    funcMap := template.FuncMap{
        "multiply": func(a, b float64) float64 {
            return a * b
        },
    }
    
    // Try different possible paths
    possiblePaths := []string{
        filepath.Join("templates", templateName),
        filepath.Join("backend", "templates", templateName),
        templateName,
    }
    
    var tmpl *template.Template
    var err error
    var errors []string
    
    for _, path := range possiblePaths {
        // Check if file exists before trying to parse it
        if _, statErr := os.Stat(path); statErr == nil {
            tmpl, err = template.New(filepath.Base(path)).Funcs(funcMap).ParseFiles(path)
            if err == nil {
                return tmpl, nil
            }
            errors = append(errors, fmt.Sprintf("Path %s exists but parse error: %v", path, err))
        } else {
            errors = append(errors, fmt.Sprintf("Path %s does not exist: %v", path, statErr))
        }
    }
    
    // Return detailed error if all paths failed
    return nil, fmt.Errorf("failed to load template from any path: %v", errors)
}
