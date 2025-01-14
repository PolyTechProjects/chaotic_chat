package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/PolyTechProjects/chaotic_chat/chat/src/config"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthGRPCClient struct {
	auth.AuthClient
}

func NewAuthClient(cfg *config.Config) *AuthGRPCClient {
	connectionUrl := fmt.Sprintf("%s:%s", cfg.Auth.AuthHost, cfg.Auth.AuthPort)
	conn, err := grpc.NewClient(connectionUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	slog.Info("Connected to Auth")
	slog.Info(connectionUrl)
	return &AuthGRPCClient{auth.NewAuthClient(conn)}
}

func (authClient *AuthGRPCClient) PerformAuthorize(ctx context.Context, r *http.Request) (*auth.AuthorizeResponse, error) {
	var accessToken string
	if r == nil {
		accessToken = metadata.ValueFromIncomingContext(ctx, "authorization")[0]
	} else {
		ctx = r.Context()
		authHeader := r.Header.Get("Authorization")
		accessToken = strings.Split(authHeader, " ")[1]
	}
	return authClient.Authorize(ctx, &auth.AuthorizeRequest{AccessToken: accessToken})
}
