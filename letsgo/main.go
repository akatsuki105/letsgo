package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	modname := flag.Arg(0)
	if modname == "" {
		fmt.Println("letsgo <name>")
		return 0
	}

	if err := github(); err != nil {
		fmt.Printf(".github/workflows/go.yml: %s", err)
		return 1
	}

	if err := vscode(); err != nil {
		fmt.Printf(".vscode/settings.json: %s", err)
		return 1
	}

	if err := gomod(modname); err != nil {
		fmt.Printf("go.mod: %s", err)
		return 1
	}

	if err := cmd(); err != nil {
		fmt.Printf("cmd/main.go: %s", err)
		return 1
	}

	if err := pkg(); err != nil {
		fmt.Printf("pkg/: %s", err)
		return 1
	}

	if err := makefile(modname); err != nil {
		fmt.Printf("makefile: %s", err)
		return 1
	}

	if err := gitignore(modname); err != nil {
		fmt.Printf(".gitignore: %s", err)
		return 1
	}

	return 0
}

func github() error {
	// mkdir .github
	if err := os.Mkdir(".github", 0777); err != nil {
		return err
	}

	// mkdir workflows
	if err := os.Mkdir(".github/workflows", 0777); err != nil {
		return err
	}

	// touch go.yml
	text := `name: Go
on:
  push:
    branches:
      - master
    paths-ignore:
      - "README.md"

jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - uses: actions/setup-go@v1
        with:
          go-version: 1.13
      - run: go build ./cmd/main.go`
	if err := ioutil.WriteFile(".github/workflows/go.yml", []byte(text), 0777); err != nil {
		return err
	}

	return nil
}

func vscode() error {
	// mkdir .vscode
	if err := os.Mkdir(".vscode", 0777); err != nil {
		return err
	}

	// touch settings.json
	text := `{
    "go.useLanguageServer": true,
    "go.alternateTools": {
        "go-langserver": "gopls"
    },
    "[go]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true
        }
    },
    "go.autocompleteUnimportedPackages": true,
    "gopls": {
        "usePlaceholders": true
    }
}`
	if err := ioutil.WriteFile(".vscode/settings.json", []byte(text), 0777); err != nil {
		return err
	}

	return nil
}

func gomod(name string) error {
	// touch go.mod
	text := `module %s

go 1.13`
	text = fmt.Sprintf(text, name)
	if err := ioutil.WriteFile("go.mod", []byte(text), 0777); err != nil {
		return err
	}

	return nil
}

func cmd() error {
	// mkdir cmd
	if err := os.Mkdir("cmd", 0777); err != nil {
		return err
	}

	// touch main.go
	text := `package main

func main() {}
`
	if err := ioutil.WriteFile("cmd/main.go", []byte(text), 0777); err != nil {
		return err
	}

	return nil
}

func pkg() error {
	if err := os.Mkdir("pkg", 0777); err != nil {
		return err
	}

	return nil
}

func makefile(name string) error {
	// touch makefile
	text := `ifdef COMSPEC
	EXE_EXT := .exe
else
	EXE_EXT := 
endif

.PHONY: build
build:
	go build -o %s$(EXE_EXT) ./cmd/main.go`
	text = fmt.Sprintf(text, name)
	if err := ioutil.WriteFile("makefile", []byte(text), 0777); err != nil {
		return err
	}

	return nil
}

func gitignore(name string) error {
	text := `%s.exe
%s`
	text = fmt.Sprintf(text, name, name)
	if err := ioutil.WriteFile(".gitignore", []byte(text), 0777); err != nil {
		return err
	}

	return nil
}
