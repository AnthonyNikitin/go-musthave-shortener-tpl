package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v10"
	"net/url"
	"strings"
)

type Configuration struct {
	Address         string
	BaseResponseURL string
}

type EnvConfiguration struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Address:         "localhost:8080",
		BaseResponseURL: "http://localhost:8080/",
	}
}

func (c *Configuration) ParseConfiguration() {

	httpPrefix := "http://"
	httpsPrefix := "https://"

	flag.Func("a", "the address where the application will start", func(s string) error {
		_, err := url.ParseRequestURI(s)
		if err != nil {
			return err
		}

		c.Address = cutPrefixes(s, httpPrefix, httpsPrefix)

		return nil
	})

	flag.Func("b", "base url address of the response", func(s string) error {
		_, err := url.ParseRequestURI(s)
		if err != nil {
			return err
		}

		c.BaseResponseURL = setSuffixes(s, "/")

		return nil
	})

	flag.Parse()

	cfg := EnvConfiguration{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	if len(cfg.ServerAddress) > 0 {
		c.Address = cutPrefixes(cfg.ServerAddress, httpPrefix, httpsPrefix)
	}

	if len(cfg.BaseURL) > 0 {
		c.BaseResponseURL = setSuffixes(cfg.BaseURL, "/")
	}
}

func cutPrefixes(s string, prefixes ...string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return s[len(prefix):]
		}
	}

	return s
}

func setSuffixes(s string, suffixes ...string) string {
	for _, suffix := range suffixes {
		if !strings.HasSuffix(s, suffix) {
			return s + "/"
		}
	}

	return s
}
