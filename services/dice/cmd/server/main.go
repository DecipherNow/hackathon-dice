package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"google.golang.org/grpc"

	"github.com/lucasmoten/project-2502/services/dice/cmd/server/config"
	"github.com/lucasmoten/project-2502/services/dice/cmd/server/methods"

	// we don't use this directly, but need it in vendor for gateway grpc plugin
	_ "github.com/ghodss/yaml"
	_ "github.com/golang/glog"
	_ "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func main() {
	var tlsServerConf *tls.Config
	var err error

	logger := zerolog.New(os.Stdout).
		With().Timestamp().Str("service", "dice").Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger.Info().Msg("starting")

	ctx, cancelFunc := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Debug().Msg("initializing config")
	if err = config.Initialize(); err != nil {
		logger.Fatal().AnErr("config.Initialize()", err).Msg("")
	}

	logger.Debug().Msg("initializing database")
	db, err := sql.Open("sqlite3", "./dice.db")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error opening sqlite db")
	}
	_, err = db.Exec("create table if not exists `githubusers` (`login` VARCHAR(64) primary key, `avatar_url` VARCHAR(256) NULL, `name` VARCHAR(128) NULL, `email` VARCHAR(256) NULL, `created_at` VARCHAR(32) NULL, `updated_at` VARCHAR(32) NULL)")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error creating table")
	}
	stmt, err := db.Prepare("insert into githubusers (login, avatar_url, name, email, created_at, updated_at) values (?,?,?,?,?,?) on conflict(login) do update set login=login where 1=0")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error preparing statement")
	}
	_, err = stmt.Exec("__cachestate__", "https://lcsc.academyofmine.com/wp-content/uploads/2017/06/Test-Logo.svg.png", "Test User", "test@deciphernow.com", "2018-12-13T00:00:00Z", "2018-12-13T00:00:00Z")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error inserting test record")
	}
	db.Close()

	if tlsServerConf, err = buildServerTLSConfigIfNeeded(logger); err != nil {
		logger.Fatal().AnErr("buildServerTLSConfigIfNeeded", err).Msg("")
	}

	logger.Debug().Str("host", viper.GetString("grpc_server_host")).
		Int("port", viper.GetInt("grpc_server_port")).
		Msg("creating listener")

	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(
			"%s:%d",
			viper.GetString("grpc_server_host"),
			viper.GetInt("grpc_server_port"),
		),
	)
	if err != nil {
		logger.Fatal().AnErr("net.Listen", err).Msg("")
	}

	opts := getTLSOptsIfNeeded(tlsServerConf)
	grpcServer := grpc.NewServer(opts...)

	methods.CreateAndRegisterServer(logger, grpcServer)

	logger.Debug().Msg("starting grpc server")
	go grpcServer.Serve(lis)

	if viper.GetBool("use_gateway_proxy") {
		logger.Debug().Msg("starting gateway proxy")
		if err = startGatewayProxy(ctx, logger); err != nil {
			logger.Fatal().AnErr("startGatewayProxy", err).Msg("")
		}
	}

	logger.Debug().Str("github_api", viper.GetString("github_api")).Str("github_token", viper.GetString("github_token")).Msg("github config")

	s := <-sigChan
	logger.Info().Str("signal", s.String()).Msg("shutting down")
	cancelFunc()
	grpcServer.Stop()
}
