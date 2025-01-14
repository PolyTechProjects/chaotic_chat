package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/config"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthGRPCClient struct {
	authClient auth.AuthClient
}

func New(cfg *config.Config) *AuthGRPCClient {
	connectionUrl := fmt.Sprintf("%s:%s", cfg.Auth.AuthHost, cfg.Auth.AuthPort)
	conn, err := grpc.NewClient(connectionUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	return &AuthGRPCClient{authClient: auth.NewAuthClient(conn)}
}

func (c *AuthGRPCClient) PerformAuthorize(ctx context.Context, r *http.Request) (*auth.AuthorizeResponse, error) {
	var accessToken string
	if r == nil {
		accessToken = metadata.ValueFromIncomingContext(ctx, "authorization")[0]
	} else {
		ctx = r.Context()
		authHeader := r.Header.Get("Authorization")
		accessToken = strings.Split(authHeader, " ")[1]
	}
	return c.authClient.Authorize(ctx, &auth.AuthorizeRequest{AccessToken: accessToken})
}
