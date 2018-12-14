package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	_ "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"net"
	"os"
	"os/signal"
	"syscall"

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
		logger.Fatal().AnErr("sqlite init", err).Msg("error creating table for githubusers")
	}
	_, err = db.Exec("create table if not exists `repositories` (`id` int  primary key, `node_id` VARCHAR(256) NULL, `name` VARCHAR(128) NULL, " +
		"`full_name` VARCHAR(256) NULL, `description` VARCHAR(256) NULL, `language` VARCHAR(32) NULL, `default_branch` VARCHAR(32) NULL, `created_at` TIMESTAMP NULL, `pushed_at` TIMESTAMP NULL, `updated_at` TIMESTAMP NULL," +
		"`fork` bool NULL,`private` bool NULL,`archived` bool NULL, `forks_count` int NULL, `network_count` int NULL, `open_issues_count` int NULL, " +
		"`stargazers_count` int NULL, `subscribers_count` int NULL, `watchers_count` int NULL, `size` int NULL)")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error creating table for repositories")
	}
	stmtGithubUsers, err := db.Prepare("insert into githubusers (login, avatar_url, name, email, created_at, updated_at) values (?,?,?,?,?,?) on conflict(login) do update set login=login where 1=0")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error preparing statement to insert githubusers cachestate")
	}
	_, err = stmtGithubUsers.Exec("__cachestate__", "https://lcsc.academyofmine.com/wp-content/uploads/2017/06/Test-Logo.svg.png", "Test User", "test@deciphernow.com", "2018-12-13T00:00:00Z", "2018-12-13T00:00:00Z")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error inserting test record for githubusers")
	}
	_, err = db.Exec("create table if not exists `githubevents` (`id` varchar(32) primary key, created_at VARCHAR(32), actor VARCHAR(64), repo_name VARCHAR(128), type VARCHAR(32), action VARCHAR(32) null, payload_comment_body TEXT null, payload_issue_number INT null, payload_issue_title TEXT null, payload_issue_body TEXT null, payload_pr_number INT null, payload_pr_title TEXT null, payload_pr_body TEXT null, payload_pr_merged_dt VARCHAR(32) null, payload_pr_merged VARCHAR(5) null, ref VARCHAR(128) null, ref_type VARCHAR(32) null, payload_size INT null, summaryline TEXT null)")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error creating table for githubevents")
	}
	stmtGithubEvents, err := db.Prepare("insert into githubevents (id, created_at, actor, repo_name, type) values (?, ?, ?, ?, ?) on conflict(id) do update set id=id where 1=0")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error preparing statement to insert githubevents cachestate")
	}
	_, err = stmtGithubEvents.Exec("0", "2018-12-13T00:00:00Z", "__cachestate__", "*/*", "CacheEvent")
	if err != nil {
		logger.Fatal().AnErr("sqlite init", err).Msg("error inserting test record for githubusers")
	}

	//stmt, err := db.Prepare("insert into githubusers (login, avatar_url, name, email, created_at, updated_at) values (?,?,?,?,?,?) on conflict(login) do update set login=login where 1=0")
	//if err != nil {
	//	logger.Fatal().AnErr("sqlite init", err).Msg("error preparing statement")
	//}
	//
	//_, err = stmt.Exec("__cachestate__", "https://lcsc.academyofmine.com/wp-content/uploads/2017/06/Test-Logo.svg.png", "Test User", "test@deciphernow.com", "2018-12-13T00:00:00Z", "2018-12-13T00:00:00Z")
	//if err != nil {
	//	logger.Fatal().AnErr("sqlite init", err).Msg("error inserting test record")
	//}

	//repos, err := db.Prepare("insert into repositories (id, node_id, name," +
	//	" full_name, description, created_at, pushed_at, updated_at, fork," +
	//	" private, archived, forks_count, network_count, open_issues_count," +
	//	" stargazers_count, subscribers_count, watchers_count, size) values" +
	//	" (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) on conflict(id) do update set id=id where 1=0")
	//if err != nil {
	//	logger.Fatal().AnErr("sqlite init", err).Msg("error preparing statement")
	//}
	//
	//_, err = repos.Exec("0", "TEST" , "__cachestate__", "TEST", "TEST", "2018-12-13T00:00:00Z", "2018-12-13T00:00:00Z", "2018-12-13T00:00:00Z", false, false, false, 1, 1, 1, 1, 1, 1, 1)
	//if err != nil {
	//	logger.Fatal().AnErr("sqlite init", err).Msg("error inserting test record")
	//}

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
