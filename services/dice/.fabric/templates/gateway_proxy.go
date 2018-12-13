package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/deciphernow/gm-fabric-go/middleware"
	"github.com/deciphernow/gm-fabric-go/tlsutil"

	pb "{{.PBImport}}"
)

func startGatewayProxy(ctx context.Context, logger zerolog.Logger) error {
	var listener net.Listener
	var err error

	mux := runtime.NewServeMux()
	var handler http.Handler = mux

	m := []middleware.Middleware{
		middleware.MiddlewareFunc(hlog.NewHandler(logger)),
	}

	if viper.GetBool("verbose_logging") {
		logger.Level(zerolog.DebugLevel)
	}

	stack := middleware.Chain(m...)
	handler = stack.Wrap(handler)

	if err = registerClient(ctx, logger, mux); err != nil {
		return errors.Wrap(err, "registerClient")
	}

	proxyAddress := fmt.Sprintf(
		"%s:%d",
		viper.GetString("gateway_proxy_host"),
		viper.GetInt("gateway_proxy_port"),
	)
	if viper.GetBool("gateway_use_tls") {
		var tlsServerConf *tls.Config

		tlsServerConf, err = tlsutil.BuildServerTLSConfig(
			viper.GetString("ca_cert_path"),
			viper.GetString("server_cert_path"),
			viper.GetString("server_key_path"),
		)
		if err != nil {
			return errors.Wrap(err, "tlsutil.BuildServerTLSConfig")
		}
		listener, err = tls.Listen("tcp", proxyAddress, tlsServerConf)
		if err != nil {
			return errors.Wrap(err, "tls.Listen failed")
		}
	} else {
		listener, err = net.Listen("tcp", proxyAddress)
		if err != nil {
			return errors.Wrap(err, "tls.Listen failed")
		}
	}

	logger.Debug().Str("service", "{{.ServiceName}}").
		Str("host", viper.GetString("gateway_proxy_host")).
		Int("port", viper.GetInt("gateway_proxy_port")).
		Msg("starting gateway proxy server")

	go http.Serve(listener, handler)

	return nil
}

func registerClient(
	ctx context.Context,
	logger zerolog.Logger,
	mux *runtime.ServeMux,
) error {
	var err error

	var clientOpts []grpc.DialOption
	if viper.GetBool("grpc_use_tls") {
		var creds credentials.TransportCredentials
		var tlsClientConf *tls.Config

		tlsClientConf, err = tlsutil.NewTLSClientConfig(
			viper.GetString("ca_cert_path"),
			viper.GetString("server_cert_path"),
			viper.GetString("server_key_path"),
			viper.GetString("server_cert_name"),
		)
		if err != nil {
			return errors.Wrap(err, "tlsutil.NewTLSClientConfig")
		}

		creds = credentials.NewTLS(tlsClientConf)
		clientOpts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		clientOpts = []grpc.DialOption{grpc.WithInsecure()}
	}

	err = pb.Register{{.ServerInterfaceName}}HandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(
			"%s:%d",
			viper.GetString("grpc_server_host"),
			viper.GetInt("grpc_server_port"),
		),
		clientOpts,
	)
	if err != nil {
		return errors.Wrap(err, "pb.Register{{.ServerInterfaceName}}HandlerFromEndpoint")
	}

	return nil
}
