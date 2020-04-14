package secretmanager

import (
	"context"
	"fmt"
	"strings"

	gcpsm "cloud.google.com/go/secretmanager/apiv1"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	pb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Client は GCP の Secret Manager にアクセスするクライアント
type Client struct {
	service *gcpsm.Client
}

// NewClient (､´･ω･)▄︻┻┳═一
func NewClient(ctx context.Context) (*Client, error) {
	service, err := gcpsm.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize GCP Secret Manager client")
	}
	return &Client{service: service}, nil
}

// ListSecretsRequest は projects.secrets.list へのリクエスト
type ListSecretsRequest struct {
	ProjectName string
	Labels      map[string]string
}

// ResourceName はリソースネームを返す
func (req *ListSecretsRequest) ResourceName() string {
	return fmt.Sprintf("projects/%s", req.ProjectName)
}

// ListSecretsResponse は projects.secrets.list のレスポンス
type ListSecretsResponse struct {
	Keys []string
}

// ListSecrets はAPIにリクエストし、シークレット一覧を取得して返す
func (cli *Client) ListSecrets(ctx context.Context, req *ListSecretsRequest) (*ListSecretsResponse, error) {
	it := cli.service.ListSecrets(ctx, &pb.ListSecretsRequest{
		Parent: req.ResourceName(),
	})
	keys := []string{}
	for {
		res, err := it.Next()
		if err != nil {
			if err != iterator.Done {
				return nil, errors.Wrap(err, "failed to list secret")
			}
			break
		}
		sps := strings.Split(res.Name, "/")
		keys = append(keys, sps[len(sps)-1])
	}
	return &ListSecretsResponse{
		Keys: keys,
	}, nil
}

// AccessSecretVersion はAPIにリクエストし、シークレット情報を取得して返す
func (cli *Client) AccessSecretVersion(ctx context.Context, req *AccessSecretVersionRequest) (*AccessSecretVersionResponse, error) {
	res, err := cli.service.AccessSecretVersion(ctx, &pb.AccessSecretVersionRequest{
		Name: req.ResourceName(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to request access sercret version")
	}
	return &AccessSecretVersionResponse{
		Value: string(res.Payload.Data),
	}, nil
}

// AccessSecretVersionRequest は projects.secrets.versions へのリクエスト
type AccessSecretVersionRequest struct {
	ProjectName string
	Key         string
	Version     string
}

// ResourceName はリソースネームを返す
func (req *AccessSecretVersionRequest) ResourceName() string {
	v := req.Version
	if v == "" {
		v = "latest"
	}
	return fmt.Sprintf("projects/%s/secret/%s/versions/%s", req.ProjectName, req.Key, v)
}

// AccessSecretVersionResponse は projects.secrets.versions のレスポンス
type AccessSecretVersionResponse struct {
	Value string
}
