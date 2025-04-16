package suite

import (
	"auth/pkg/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg     *config.Config
	Client  *http.Client
	BaseURL string
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.ServerConfig.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	return ctx, &Suite{
		T:       t,
		Cfg:     cfg,
		BaseURL: fmt.Sprintf("http://localhost:%d/api/v1/auth", cfg.ServerConfig.Port),
		Client:  &http.Client{Timeout: cfg.ServerConfig.Timeout},
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/local.yaml"
}
