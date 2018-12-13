package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/deciphernow/gm-fabric-go/tlsutil"
)

func main() {
	os.Exit(run())
}

func run() int {
	var testCertDir string
	var uri string
	var err error

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr})

	pflag.StringVar(
		&uri,
		"uri",
		"",
		"URI to send to server",
	)
	pflag.StringVar(
		&testCertDir,
		"test-cert-dir",
		"",
		"(if TLS) directory holding test certificates",
	)
	pflag.Parse()
	if uri == "" {
		logger.Error().Msg("You must supply a URI. (--uri)")
		return 1
	}

	logger.Info().Str("http-client", "dice").
		Str("URI", uri).
		Str("test-certs", testCertDir).Msg("starting")

	client, err := newClient(testCertDir)
	if err != nil {
		logger.Error().AnErr("newClient", err).Msg("")
		return 1
	}

	if err = runURITest(client, uri); err != nil {
		logger.Error().AnErr("runGatewayTest", err).Msg("")
		return 1
	}

	logger.Info().Str("http client", "dice").Msg("terminating normally")
	return 0
}

func newClient(testCertDir string) (*http.Client, error) {
	var client http.Client
	var err error

	if testCertDir != "" {
		var tlsConf *tls.Config

		tlsConf, err = tlsutil.NewTLSClientConfig(
			filepath.Join(testCertDir, "intermediate.crt"),
			filepath.Join(testCertDir, "localhost.crt"),
			filepath.Join(testCertDir, "localhost.key"),
			"localhost",
		)
		if err != nil {
			return nil, errors.Wrap(err, "tlsutil.NewTLSClientConfig")
		}

		transport := http.Transport{
			TLSClientConfig: tlsConf,
		}

		client.Transport = &transport
	}

	return &client, nil
}

func runURITest(client *http.Client, uri string) error {
	httpResp, err := client.Get(uri)
	if err != nil {
		return errors.Wrap(err, "client.Get")
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return errors.Errorf(
			"invalid HTTP status (%d) %s",
			httpResp.StatusCode,
			httpResp.Status,
		)
	}

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll")
	}

	fmt.Println()
	fmt.Println("output of %s", uri)
	fmt.Println(string(respBytes))
	fmt.Println()

	return nil
}
