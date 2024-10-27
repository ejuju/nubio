# Setup guide for HTTPS server with Systemd on Debian (with CI/CD)

TODO

---

## Add Github workflow to your repo

In order to deploy our changes, we need to:
- Copy the new `config.json` to `/var/www/config.json`
- Restart the server: `systemctl restart website`

This can be achieved with a Github workflow (`.github/workflows/deploy.yaml`):
```yaml
TODO
```

## Setup the remote working directory

SSH into your development server

Create new directories:
```bash
echo "Setting up working directory..." && \
mkdir /var/www && \
mkdir /var/www/tls
```

Add your config (in `/var/www/config.json`):
```json
{
    "domain": "alexdoe.example",
    "tls_dirpath": "tls",
    "tls_email_addr": "contact@alexdoe.example",
    "resume_path": "resume.json"
}
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
