package scalar

import (
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/caiolandgraf/grove-base/internal/config"
)

var customCSS = `
	.scalar-app {
    --scalar-font: 'Inter', sans-serif;

    --scalar-background-1: #0D0D0D;
    --scalar-background-2: #181818;
    --scalar-background-3:  #0A0A0A;
    --scalar-background-accent: #232323;

    --scalar-border-color: #282828;

    --scalar-color-1: #c2c2c2;
    --scalar-color-accent: #F3A8F7;
    --scalar-color-3: #E5B28A;
  }
`

func NewUI(specURL string) http.Handler {
	baseURL := config.Env.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	fullURL := baseURL + specURL

	html, err := scalargo.NewV2(
		scalargo.WithSpecURL(fullURL),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle(config.Env.AppName),
			scalargo.WithKeyValue(
				"description",
				config.Env.AppDesc,
			),
			scalargo.WithKeyValue(
				"ogDescription",
				config.Env.AppOGDC,
			),
		),
		scalargo.WithHideDarkModeToggle(),
		scalargo.WithHideDownloadButton(),
		scalargo.WithPersistAuth(true),
		scalargo.WithOverrideCSS(customCSS),
	)
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	})
}
