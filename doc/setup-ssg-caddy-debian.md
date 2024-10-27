# Setup guide for SSG with Caddy on Debian (with CI/CD)

In this guide, we will setup a self-hosted online resume
using Caddy, Debian, and Github Actions.

## Install Go on your development machine

> Skip this part if Go is already installed on your machine.

To install Go 1.23.2 on a Linux/amd64 OS:
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

Note: You may need to change the Go version and CPU architecture in the commands above.

## Install (or update) Nubio on your local development machine

> Skip this part if you already have the lastest Nubio version installed on your machine.

To install Nubio, run:
```bash
go install github.com/ejuju/nubio@latest
```

## Setup your Git repo for local development

Create your resume configuration in `/resume.json`.
A sample configuration is available here: [sample resume.json](/resume.json)

We will now setup a local development server to generate 
and view our exports in a web browser:

Create your local `/server.json`:
```json
{
    "address": ":8080",
    "resume_path": "resume.json"
}
```

Note: You may need to change the port specified in the `address` field
depending on your local setup and preferences.

Now let's check your config files and start the local server:
```bash
echo "Checking config files..." && \
nubio check-resume-config resume.json && \
nubio check-server-config server.json && \
echo "Starting server..." && \
nubio run
```

Your website is now running on [localhost:8080](http://localhost:8080).
Use CTRL+C to stop the server.

## Get a host with a publically accessible IP address.

Here, we will be using a Debian 12 VPS.
You may use the provider of your choice for this.

<!-- TODO: Create and use non-root user -->
Warning: We will be using the root user here,
for a more secure setup, you should create a dedicated user on the server host.

## Update your DNS

- Ensure your domain resolves to your server's IP address(es).
- Ensure you have a CNAME on `www.` pointing to `alexdoe.example.`.

## Generate static website

To generate your website files, run:
```bash
nubio ssg config.json .out
```

## Setup Caddy on the remote server

```bash
ssh root@alexdoe.example
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
mkdir -p /var/www/alexdoe
```

Setup your Caddy file in `/etc/caddy/Caddyfile`:
```
alexdoe.example {
	root * /var/www/alexdoe
	file_server
}

www.alexdoe.example {
	redir https://alexdoe.example{uri} permanent
}
```

Restart Caddy:
```bash
systemctl restart caddy && systemctl status caddy
```

## Copy static files to remote server

From your local machine, copy the static files created by `nubio` to the remote server:

```bash
scp -r .out/* root@alexdoe.example:/var/www/alexdoe
```

Your website is now up and running!

## Setup CI/CD

On your local machine, generate a new SSH key pair: `ssh-keygen`.

Copy the public key on the remote server (append to `authorized_keys file`):
```bash
ssh root@alexdoe.example
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
          go install github.com/ejuju/nubio@v0.5.0
          mkdir .out
          nubio check-resume-config resume.json
          nubio ssg resume.json .out
          echo "$SSH_KEY" > "$SSH_KEY_PATH"
          chmod 0600 "$SSH_KEY_PATH"
          scp \
              -i "$SSH_KEY_PATH" \
              -o StrictHostKeyChecking=no \
              -o UserKnownHostsFile=/dev/null \
              -r .out/* "$SSH_USERNAME"@alexdoe.example:/var/www/alexdoe
```

At this point this is what your repositery should contain:
- `/.github/workflows/deploy.yaml`: Github action for continuous deployment.
- `/resume.json`: Resume configuration.
- `/server.json`: Used for local development server.

Push the code:
```bash
git add .
git commit -m "cicd: add cicd"
git push
```

CI/CD is setup!
You can now make modifications to your `resume.json` and
the new changes will be deployed to your server when you push.
