# Example setup for SSG with Caddy on Debian

Done on a fresh Debian 11 install.

We use the root user here, but in reality you should create a dedicated user on the server host.

You should also setup CICD to deploy the new static files on code change:
SSG would be done in an automated pipeline: where you would install Nubio 
and generate website then push to the actual server.

## Install Caddy

Install and start the Caddy daemon:
```bash
echo "Installing Caddy..." && \
    sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl && \
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg && \
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list && \
    sudo apt update && \
    sudo apt install caddy && \
    echo "Done!"
```

Source: https://caddyserver.com/docs/install#debian-ubuntu-raspbian

Ensure Caddy is running:
```bash
systemctl status caddy
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

## Generate static website

```bash
# Add your profile.json
cd /tmp
vim profile.json

mkdir -p /var/www/mysite # Create working directory
nubio ssg profile.json /var/www/mysite # Generate website files
```

## Update your DNS

- Ensure your domain resolves to your server's IP address(es).
- Ensure you have a CNAME on `www.` pointing to `www.mysite.example.`.

## Setup Caddy

Setup your Caddy file in `/etc/caddy/Caddyfile`:
```
mysite.example {
	root * /var/www/mysite
	file_server
}

www.mysite.example {
	redir https://mysite.example{uri} permanent
}
```

## Restart Caddy
Restart Caddy:
```bash
systemctl restart caddy && systemctl status caddy
```

That's it! You can now visit your website.
