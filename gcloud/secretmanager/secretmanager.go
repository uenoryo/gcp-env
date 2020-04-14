package secretmanager

import (
	"context"
	"fmt"

	gcpsm "cloud.google.com/go/secretmanager/apiv1"
	"github.com/pkg/errors"
	pb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Client は GCP の Secret Manager にアクセスするクライアント
type Client struct {
	service *gcpsm.Client
}

// NewClient (､´･ω･)▄︻┻┳═一
func NewClient() (*Client, error) {
	service, err := gcpsm.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize GCP Secret Manager client")
	}
	return &Cleint{service: service}
}

// AccessSecretVersion はAPIにリクエストし、シークレット情報を取得して返す
func (cli *Client) AccessSecretVersion(ctx context.Context, req *AccessSecretVersionRequest) (*AccessSecretVersionResponse, error) {
	res, err := cli.service.AccessSecretVersion(ctx, pb.AccessSecretVersionRequest{
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

// ResourceName リソースネームを返す
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
