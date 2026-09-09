package view

var toolFileIcons = map[string]string{
	".dockerignore":    "docker",
	"build.gradle":     "java",
	"bun.lock":         "lock",
	"bun.lockb":        "lock",
	"cargo.lock":       "lock",
	"cargo.toml":       "rust",
	"composer.json":    "php",
	"composer.lock":    "lock",
	"containerfile":    "docker",
	"gemfile":          "ruby",
	"gemfile.lock":     "lock",
	"go.mod":           "go",
	"go.sum":           "go",
	"go.work":          "go",
	"gradlew":          "terminal",
	"makefile":         "makefile",
	"package.json":     "javascript",
	"pipfile":          "python",
	"pipfile.lock":     "lock",
	"pom.xml":          "java",
	"pyproject.toml":   "python",
	"rakefile":         "ruby",
	"requirements.txt": "python",
	"settings.gradle":  "java",
}

var javascriptCompoundIcons = map[string]string{
	".spec.cjs": "javascript",
	".spec.js":  "javascript",
	".spec.mjs": "javascript",
	".test.cjs": "javascript",
	".test.js":  "javascript",
	".test.mjs": "javascript",
}

var typescriptCompoundIcons = map[string]string{
	".spec.ts":  "typescript",
	".spec.tsx": "typescript",
	".test.ts":  "typescript",
	".test.tsx": "typescript",
	".tsx":      "typescript",
}
