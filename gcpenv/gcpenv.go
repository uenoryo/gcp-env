package gcpenv

import (
	"context"
	"fmt"

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
	res, err := client.ListSecrets(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to list secrets")
	}

	for _, key := range res.Keys {
		req := &secretmanager.AccessSecretVersionRequest{
			ProjectName: "chitoi",
			Key:         key,
			//			Version:     "1",
		}
		res, err := client.AccessSecretVersion(ctx, req)
		if err != nil {
			return errors.Wrapf(err, "failed to access secret. key:[%s]", key)
		}
		fmt.Printf("%s=%s\n", key, res.Value)
	}
	return nil
}

// Config (､´･ω･)▄︻┻┳═一
type Config struct{}
