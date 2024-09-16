package main

import (
	"flag"
	"net/url"
	"os"
)

// AppConfig contains all the application configuration directives.
type AppConfig struct {
	// The base URL to use for all shortened URLs.
	baseURL *url.URL
	// The path where the BoltDB database file should be stored.
	dbFile string
	// Whether HTTPS should be enabled or not.
	// This is automatically determined depending on whether httpsCertFile and httpsKeyFile are passed.
	httpsEnabled  bool
	httpsCertFile string
	httpsKeyFile  string
	// The listen address to use for the HTTP server. This is automatically determined if not explicitly specified.
	listenAddress string
	logLevel      string
}

// getAppConfiguration will parse the environment, command-line flags and application defaults to the final application configuration.
func getAppConfiguration() (*AppConfig, error) {
	var err error
	// Set the base URL to a string for now, so I only have to parse it once
	var rawBaseURL string

	c := &AppConfig{}

	prefix := "URL_SHORTENER_"
	rawBaseURL = os.Getenv(prefix + "BASE_URL")
	c.dbFile = os.Getenv(prefix + "DB_FILE")
	c.httpsCertFile = os.Getenv(prefix + "HTTPS_CERT_FILE")
	c.httpsKeyFile = os.Getenv(prefix + "HTTPS_KEY_FILE")
	c.listenAddress = os.Getenv(prefix + "LISTEN_ADDRESS")
	c.logLevel = os.Getenv(prefix + "LOG_LEVEL")

	flag.StringVar(&rawBaseURL, "base-url", rawBaseURL, "the base URL to use for all shortened URLs")
	flag.StringVar(&c.dbFile, "db-file", c.dbFile, "the path where the BoltDB database file should be stored, defaults to 'url-shortener.db' in the current working directory")
	flag.StringVar(&c.httpsCertFile, "https-cert-file", c.httpsCertFile, "the TLS certificate file to use for HTTPS")
	flag.StringVar(&c.httpsKeyFile, "https-key-file", c.httpsKeyFile, "the TLS certificate key file to use for HTTPS")
	flag.StringVar(&c.listenAddress, "listen-address", c.listenAddress, "the address that the application should listen on")
	flag.StringVar(&c.logLevel, "log-level", c.logLevel, "the log level to use, e.g. DEBUG, INFO (default), WARN or ERROR")
	flag.Parse()

	if rawBaseURL != "" {
		c.baseURL, err = url.ParseRequestURI(rawBaseURL)
		if err != nil {
			return c, err
		}
	} else {
		// We will just try to determine this from the request itself if it is not set, so don't error
		c.baseURL = nil
	}

	if c.dbFile == "" {
		// If no database file is explicitly given, default to 'url-shortener.db' in the current working directory
		c.dbFile = "url-shortener.db"
	}

	// If both a certificate and corresponding key file path have been passed, HTTPS will be enabled
	if c.httpsCertFile != "" && c.httpsKeyFile != "" {
		c.httpsEnabled = true
	} else {
		c.httpsEnabled = false
	}

	// If no listen address has been passed, set an appropriate default depending on whether HTTPS has been enabled
	if c.listenAddress == "" {
		if c.httpsEnabled {
			c.listenAddress = ":443"
		} else {
			c.listenAddress = ":80"
		}
	}

	if c.logLevel == "" {
		c.logLevel = "INFO"
	}

	return c, nil
}
