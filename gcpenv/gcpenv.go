package gcpenv

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/uenoryo/gcp-env/gcloud/secretmanager"
)

// GCPEnv は環境変数を扱う
// GOOGLE_APPLICATION_CREDENTIALS に認証情報が必要です
type GCPEnv struct {
	config *Config
}

// New (､´･ω･)▄︻┻┳═一
func New(conf *Config) *GCPEnv {
	return &GCPEnv{config: conf}
}

// Fetch は環境変数を取得する
func (env *GCPEnv) Fetch(ctx context.Context) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize client")
	}

	req := &secretmanager.ListSecretsRequest{
		ProjectName: "chitoi",
	}
	secrets, err := client.ListSecrets(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to list secrets")
	}
	log.Println(secrets)
	return nil
}

// Config (､´･ω･)▄︻┻┳═一
type Config struct{}
