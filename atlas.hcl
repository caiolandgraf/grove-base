data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  url = "postgres://postgres:postgres@localhost:5432/mcs_dctfweb_sender?sslmode=disable"

  dev = "docker://postgres/15/dev"

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "dev" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")

  dev = "docker://postgres/15/dev"

  migration {
    dir = "file://migrations"
  }
}

env "production" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")

  migration {
    dir = "file://migrations"
  }
}
