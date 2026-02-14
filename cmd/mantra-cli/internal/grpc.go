package internal

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/daedaluz/mantra-cli/lib/grpc/client"
	admin "github.com/daedaluz/mantra-cli/lib/grpc/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func GetGRPCClient(cmd *cobra.Command) *grpc.ClientConn {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		panic(err)
	}
	plaintext, err := cmd.Flags().GetBool("plaintext")
	if err != nil {
		panic(err)
	}
	skipVerify, err := cmd.Flags().GetBool("skip-verify")
	if err != nil {
		panic(err)
	}

	var opts []grpc.DialOption
	if plaintext {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else if skipVerify {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	} else {
		creds := credentials.NewClientTLSFromCert(nil, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	if client, err := grpc.NewClient(server, opts...); err != nil {
		panic(err)
	} else {
		return client
	}
}

func GetAdminAPIClient(cmd *cobra.Command) (context.Context, admin.PlatformAdminServiceClient) {
	cli := GetGRPCClient(cmd)
	md := metadata.Pairs(
		"x-api-key", cmd.Flag("api-key").Value.String(),
	)
	ctx := metadata.NewOutgoingContext(cmd.Context(), md)
	return ctx, admin.NewPlatformAdminServiceClient(cli)
}

func GetDomainAdminAPIClient(cmd *cobra.Command) (context.Context, admin.DomainAdminServiceClient) {
	cli := GetGRPCClient(cmd)
	md := metadata.Pairs(
		"x-client-id", cmd.Flag("client-id").Value.String(),
		"x-client-secret", cmd.Flag("client-secret").Value.String(),
		"origin", fmt.Sprintf("https://%s", cmd.Flag("domain").Value.String()),
	)
	ctx := metadata.NewOutgoingContext(cmd.Context(), md)
	return ctx, admin.NewDomainAdminServiceClient(cli)
}

func GetClientApiClient(cmd *cobra.Command) (context.Context, client.ClientServiceClient) {
	cli := GetGRPCClient(cmd)
	md := metadata.Pairs(
		"x-client-id", cmd.Flag("client-id").Value.String(),
		"x-client-secret", cmd.Flag("client-secret").Value.String(),
		"origin", fmt.Sprintf("https://%s", cmd.Flag("domain").Value.String()),
	)
	ctx := metadata.NewOutgoingContext(cmd.Context(), md)
	return ctx, client.NewClientServiceClient(cli)
}
