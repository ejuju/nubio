# Setup for SSG with Caddy on Debian (with CI/CD)

Done on a fresh Debian 11 install.

Note:
- We use the root user here, but in reality you should create a dedicated user on the server host.
- Whenever you see `mysite` or `mysite.example`, you should replace that with your actual domain name.

## Install Go on your local machine

To install Go on a Linux OS:
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

## Install Nubio on your local machine

```bash
go install github.com/ejuju/nubio@latest
```

## Generate static website

```bash
vim profile.json # Add your profile.json
nubio ssg profile.json .out # Generate website files
```

## Copy static files to remote server

```bash
scp -r .out/* root@mysite.example:/var/www/mysite
```

## Update your DNS

- Ensure your domain resolves to your server's IP address(es).
- Ensure you have a CNAME on `www.` pointing to `www.mysite.example.`.

## Setup Caddy on the remote server

```bash
ssh root@mysite.example
```

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

Create working directory (where static files will be served from):
```bash
mkdir -p /var/www/mysite # Create working directory
```

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

Restart Caddy:
```bash
systemctl restart caddy && systemctl status caddy
```

Your website is now up and running!

## Setup CI/CD with Github workflow

Setup a Git repo containing your `profile.json` file.

Publish the repo to your Github (visibility: private)

Generate a new SSH key pair: `ssh-keygen`

Copy the public key on the remote server (append to `authorized_keys file`):
```bash
ssh root@mysite.example
vim ~/.ssh/authorized_keys
```

Add secrets (in "Settings" > "Secrets and variables" > "Actions" > "New repository secret"):
- `SSH_KEY`: The SSH private key
- `SSH_USERNAME`: Your SSH username (ex: `root` or `github`)

Note: replace Go and Nubio versions with more recent ones if needed.

Add a `.github/workflows/deploy.yaml`:
```yaml
on:
  push:
    branches:
      - "main"

jobs:
  deploy:
    name: Generate and deploy website files.
    runs-on: ubuntu-latest
    env:
      SSH_KEY: ${{ secrets.SSH_KEY }}
      SSH_USERNAME: ${{ secrets.SSH_USERNAME }}
      SSH_KEY_PATH: "ssh_key"
    steps:
      - name: Git-checkout code.
        uses: actions/checkout@v3
      - name: Setup Go.
        uses: actions/setup-go@v4
        with:
          go-version: "1.21.7"
      - name: Build and deploy.
        run: |
          go install github.com/ejuju/nubio@v0.2.7
          mkdir .out
          nubio check-profile profile.json
          nubio ssg profile.json .out
          echo "$SSH_KEY" > "$SSH_KEY_PATH"
          chmod 0600 "$SSH_KEY_PATH"
          scp \
              -i "$SSH_KEY_PATH" \
              -o StrictHostKeyChecking=no \
              -o UserKnownHostsFile=/dev/null \
              -r .out/* "$SSH_USERNAME"@mysite.example:/var/www/mysite
```

Push the code:
```bash
git add .
git commit -m "cicd: add github workflow"
git push
```

CI/CD is setup!
You can now make modifications to your profile.json and
the new changes will be deployed when you push.
