# Setup guide for HTTPS server with Systemd on Debian (with CI/CD)

TODO

---

## Setup your DNS

- Add A/AAAA records pointing to your server IP address.
- Add a CNAME from `www` to your `alexdoe.example.`

## Create a new user for CICD

TODO

## Generate new SSH key pair

TODO

## Authorize CICD user

TODO

## Add Github workflow to your repo

In order to deploy our changes, we need to:
- Copy the new `config.json` to `/var/www/config.json`
- Restart the server: `systemctl restart website`

This can be achieved with a Github workflow (`.github/workflows/deploy.yaml`):
```yaml
TODO
```

## Setup your working directory

Create new directories:
```bash
mkdir /var/www
mkdir /var/www/tls
```

Add your config (in `/var/www/config.json`):
```json
{
    "domain": "alexdoe.example",
    "tls_dirpath": "tls",
    "tls_email_addr": "contact@alexdoe.example",
    // ...
}
```

## Install Go

```bash
echo "Installing Go..." && \
    rm -rf /usr/local/go ; \
    cd /tmp/ && \
    wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz && \
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile && \
    echo "export PATH=\$PATH:/root/go/bin" >> ~/.profile && \
    source ~/.profile && \
    go version
```

## Install Nubio

```bash
go install github.com/ejuju/nubio@latest
```

## Setup the new Systemd daemon

Create Systemd service file (`/etc/systemd/system/website.service`):
```conf
[Unit]
Description=Website
After=network.target

[Service]
Type=simple
Restart=always
User=root
Group=root
WorkingDirectory=/var/www/
ExecStart=/root/go/bin/nubio

[Install]
WantedBy=multi-user.target
```

```bash
echo "Setting up Systemd..." && \
    systemctl daemon-reload && \
    systemctl enable website && \
    systemctl start website && \
    systemctl status
```

The website is now available!

To inspect logs, run `journalctl -u website`
