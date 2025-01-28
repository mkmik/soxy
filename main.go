// The soxy command is a simple-stupid reverse proxy that uses Go's standard http library
// which supports out of the box using HTTP and SOCKS proxies towards the upstream.
package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime/debug"

	"github.com/alecthomas/kong"
)

// set by goreleaser
var version = "(devel)"

type Context struct {
	*CLI
}

type CLI struct {
	From             string `name:"from" required:"" help:"Source address to proxy from"`
	To               string `name:"to" required:"" help:"Destination address to proxy to"`
	ChangeHostHeader bool   `name:"change-host-header" help:"Change the Host header to the target host"`

	Version kong.VersionFlag `name:"version" help:"Print version information and quit"`
}

func (cmd *CLI) Run(cli *Context) error {
	// Parse the target URL
	targetURL, err := url.Parse(cmd.To)
	if err != nil {
		return err
	}

	// Create a reverse proxy
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			if cmd.ChangeHostHeader {
				req.Host = targetURL.Host
			}
		},
	}

	// Start the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	return http.ListenAndServe(cmd.From, nil)
}

func getVersion() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	// otherwise fallback to the version set by goreleaser
	return version
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Description(`Simple reverse proxy with CLI configuration. Supports using a proxy for the upstream with the HTTPS_PROXY=socks5://host:port env var`),
		kong.UsageOnError(),
		kong.Vars{
			"version": getVersion(),
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	err := ctx.Run(&Context{CLI: &cli})
	ctx.FatalIfErrorf(err)
}
