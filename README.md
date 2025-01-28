# soxy

## What

A simple reverse proxy that supports SOCKS5 for upstream.

It uses simple CLI args for config so it's suitable for quick oneliners.

## Why

There are many good reverse proxy servers out there but they have either one or more of the following problems:

1. They are hard to configure
2. They don't support using a SOCKS5 proxy for upstream
3. They support a proxy but require that the remote name resolves DNS locally

I wrote this tool because I had a server behind a VPN which I can access with a SOCKS5 proxy (e.g. tailscale) but, for reasons,
I cannot make my client use SOCKS5. So instead I just expose that one target HTTP(s) server over a local port and have
the reverse proxy talk to its upstream using the SOCKS proxy.

## Install

```bash
go install mkm.pub/soxy@latest
```

## Use

If you run:

```bash
soxy reverse-proxy --from=localhost:8081 --to=https://some.remote.stuff.com --change-host-header
```

Soxy will listen on port 8081 locally and forward all the requests to https://some.remote.stuff.com.

If you also pass the `--change-host-header` flag, the remote server will see itself in the `Host` header.
Otherwise the `Host` header will likely contain `localhost:8081` (details vary on how you ultimately end up using this tool).

### Upstream SOCKS

Imagine you need to connect to proxy some host that is behind a 

```bash
HTTPS_PROXY=socks5://localhost:1080 soxy reverse-proxy --from=localhost:8081 --to=https://some.internal.stuff.com --change-host-header

```



