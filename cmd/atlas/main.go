package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	_ "github.com/caiolandgraf/go-project-base/internal/modules" // load modules → model init() registrations
	"github.com/caiolandgraf/go-project-base/internal/app/database"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(database.All()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	if _, err := io.WriteString(os.Stdout, stmts); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write gorm schema: %v\n", err)
		os.Exit(1)
	}
}
