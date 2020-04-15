package gcpenv

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/uenoryo/gcp-env/gcloud/secretmanager"
)

// GCPEnv は環境変数を扱う
// GOOGLE_APPLICATION_CREDENTIALS に認証情報が必要です
type GCPEnv struct {
	config *Config
	values *sync.Map
}

// New (､´･ω･)▄︻┻┳═一
func New(conf *Config) *GCPEnv {
	return &GCPEnv{
		config: conf,
		values: &sync.Map{},
	}
}

// Fetch は環境変数を取得する
func (env *GCPEnv) Fetch(ctx context.Context) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize client")
	}

	req := &secretmanager.ListSecretsRequest{
		ProjectName: env.config.ProjectName,
	}
	res, err := client.ListSecrets(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to list secrets")
	}

	for _, key := range res.Keys {
		req := &secretmanager.AccessSecretVersionRequest{
			ProjectName: env.config.ProjectName,
			Key:         key,
			Version:     env.config.Version,
		}
		res, err := client.AccessSecretVersion(ctx, req)
		if err != nil {
			return errors.Wrapf(err, "failed to access secret. key:[%s]", key)
		}
		env.values.Store(key, res.Value)
	}
	return nil
}

// Map は取得したデータをmapで返す
func (env *GCPEnv) Map() map[string]string {
	m := map[string]string{}
	env.values.Range(func(key, value interface{}) bool {
		k, kOk := key.(string)
		v, vOk := value.(string)
		if kOk && vOk {
			m[k] = v
		}
		return true
	})
	return m
}

// Config (､´･ω･)▄︻┻┳═一
type Config struct {
	ProjectName string
	Version     string
}
