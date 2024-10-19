# Setup

## Install

Install `go install github.com/ejuju/nubio@latest`

## Start

Setup your working directory:
1. Add a `config.json` file (see `example.prod.conf.json`)
1. Add a `profile.json` file (see `example.profile.json`)
1. Optional: Add a `pgp.asc` file (PGP key file) (see `example.pgp.asc`)
1. Ensure TLS directory exists (ex: `mkdir my-certs`)

Then your choose your preferred way for setting up the daemon (RC, Docker, systemd, etc.).

Here's an example for Systemd:

Setup the service:
- For systemd, create a service file and enable/start the daemon (use `Restart=Always`).
- For container runtimes, build an image and run the container.

## Ensure that the service is accessible from the internet

```
curl https://mysite.example/ping
```

## Maintain

Update your install with: `go install github.com/ejuju/nubio@v0.0.1`
