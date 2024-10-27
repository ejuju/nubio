# Nubio: Self-hosted online resume for developers

Nubio can be used as a static site generator, CLI, HTTP(S) server or Go library.

## Introduction

### Features

- Export your resume as PDF, JSON, plain text, Markdown, or static HTML page.
- Serve your resume as a website (or generate static website files).
- Single executable.
- Auto HTTPS (get and renew certs using ACME).

### Live demo

- Website: [juliensellier.com](https://juliensellier.com/)
- PDF export: [juliensellier.com/resume.pdf](https://juliensellier.com/resume.pdf)
- JSON export: [juliensellier.com/resume.json](https://juliensellier.com/resume.json)
- Plain-text export: [juliensellier.com/resume.txt](https://juliensellier.com/resume.json)
- Markdown export: [juliensellier.com/resume.md](https://juliensellier.com/resume.json)

### 3rd-party Go dependencies

- For PDF generation: [github.com/go-pdf/fpdf](https://github.com/go-pdf/fpdf)
- For HTTPS/ACME support: [golang.org/x/crypto](https://golang.org/x/crypto)

---

## Usage

### Installation

An executable can be built from source using
and installed locally using:
```bash
go install github.com/ejuju/nubio@latest
```

### Configuration

Your server configuration and resume information is stored in a single JSON file,
usually named `config.json`.

Check out an example in [config.json](/config.json).

You can check the validity of your `config.json` with:

```bash
nubio check config.json
```

### Generating a static website (SSG)

```bash
nubio ssg config.json static/
```

For an example setup with Caddy on Debian, check out:
[/doc/setup-ssg-caddy-debian.md](/doc/setup-ssg-caddy-debian.md)

### Generating exports via CLI

You can also use Nubio as a CLI to generate static file exports,
for example, to export your config as a PDF:

```bash
nubio pdf /path/to/config.json /path/to/output.pdf
```

### Running as HTTP(S) server

First, you'll need to configure your `config.json` file with the necessary information.

Example config fields for a HTTPS server:
```json
{
    "domain": "mysite.example",
    "tls_dirpath": "tls",
    "tls_email_addr": "contact@mysite.example",
}
```

The field `tls_dirpath` indicates a directory where TLS certificates will be stored.
When using HTTPS, the server uses ports `80` and `443` by default,
there's no need to use the `address` field (it is ignored).

Example fields for a HTTP server behind a reverse proxy:
```json
{
    "domain": "mysite.example",
    "address": ":8080",
    "true_ip_header": "X-Forwarded-For",
}
```

To start the server, run:
```bash
nubio run config.json
```

Notes:
- You can also simply run `nubio run` which by default will look
  for `config.json` in its current working directory.

### Using custom CSS

In order to add custom CSS, use the corresponding config field:
- `custom_css` to import the CSS from the config file field value.
- `custom_css_path` to import the CSS from a file, when provided, overwrites `custom_css`.

Check out the [HTML template file](/pkg/nubio/resume.html.gotmpl) to see how to select
the desired elements.

Here's an example of some custom CSS to set the UI to "light mode":
```css
:root {
    --color-fg-0: hsl(0, 0%, 5%);
    --color-fg-1: hsl(0, 0%, 20%);
    --color-fg-2: hsl(0, 0%, 35%);
    --color-bg-0: hsl(0, 0%, 95%);
    --color-bg-1: hsl(0, 0%, 92%);
    --color-bg-2: hsl(0, 0%, 85%);
}
```

### Embedding in your Go program

- Export resume to PDF: `nubio.ExportPDF(w, resume)`
- Export resume to HTML: `nubio.ExportHTML(w, resume)`
- Resume type definition: see `nubio.Resume`
- And more...

Official package documentation is available here:
[pkg.go.dev/github.com/ejuju/nubio/pkg/nubio](https://pkg.go.dev/github.com/ejuju/nubio/pkg/nubio)
