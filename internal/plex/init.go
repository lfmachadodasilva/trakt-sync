package plex

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"trakt-sync/internal/config"
	"trakt-sync/internal/ctxutils"

	"github.com/LukeHagar/plexgo"
)

func InitPlex(ctx *context.Context, cfg *config.ConfigEntity) (context.Context, error) {
	if ctx == nil || cfg == nil {
		return nil, fmt.Errorf("context or config cannot be nil")
	}

	if cfg.Plex == nil || cfg.Plex.BaseURL == "" || cfg.Plex.APIKey == "" {
		return nil, fmt.Errorf("Plex BaseURL and APIKey must be set in the configuration")
	}

	// Parse the Plex BaseURL to extract protocol, IP, and port
	protocol, ip, port, err := parseURLComponents(cfg.Plex.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Determine the server protocol based on the URL scheme
	var serverProtocol plexgo.ServerProtocol
	if protocol == "http" {
		serverProtocol = plexgo.ServerProtocolHTTP
	} else if protocol == "https" {
		serverProtocol = plexgo.ServerProtocolHTTPS
	} else {
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}

	s := plexgo.New(
		plexgo.WithSecurity(cfg.Plex.APIKey),
		plexgo.WithIP(ip),
		plexgo.WithProtocol(serverProtocol),
		plexgo.WithPort(port),
	)

	// Store the Plex client in the context for later use
	ctxReturn := context.WithValue(*ctx, ctxutils.ContextPlexSdkKey, s)

	return ctxReturn, nil
}

func parseURLComponents(tmp string) (protocol string, ip string, port string, err error) {
	parsedURL, err := url.Parse(tmp)
	if err != nil {
		return "", "", "", err
	}

	protocol = parsedURL.Scheme
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) == 2 {
		ip = hostParts[0]
		port = hostParts[1]
	} else {
		ip = parsedURL.Host
		port = "32400"
	}

	return protocol, ip, port, nil
}
