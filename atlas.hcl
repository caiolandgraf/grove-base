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
  url = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}?sslmode=${getenv("DB_SSLMODE")}"

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
  url = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}?sslmode=${getenv("DB_SSLMODE")}"

  dev = "docker://postgres/15/dev"

  migration {
    dir = "file://migrations"
  }
}

env "production" {
  src = data.external_schema.gorm.url
  url = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}?sslmode=${getenv("DB_SSLMODE")}"
  migration {
    dir = "file://migrations"
  }
}
