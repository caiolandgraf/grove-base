package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/caiolandgraf/grove-base/internal/app/database"
	_ "github.com/caiolandgraf/grove-base/internal/modules" // load modules → model init() registrations
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
