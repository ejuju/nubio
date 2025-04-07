# Nubio: Self-hosted online resume tailored for developers

Nubio can be used as a static site generator, CLI, HTTP(S) server or Go library.

> Nubio is still being designed and developed. Use at your own risk.
> No retro-compatibility promises are made until v1 is reached.

## Introduction

### Features

- Configure your resume with a single JSON file.
- Export your resume as HTML, PDF, or JSON.
- Serve your resume as a website (or generate static website files).
- Auto HTTPS (get and renew certs using ACME).
- Single executable.

### Live demo

- Website: [juliensellier.com](https://juliensellier.com/)
- PDF export: [juliensellier.com/resume.pdf](https://juliensellier.com/resume.pdf)
- JSON export: [juliensellier.com/resume.json](https://juliensellier.com/resume.json)

### 3rd-party dependencies

- PDF generation: [github.com/go-pdf/fpdf](https://github.com/go-pdf/fpdf)

---

## Usage

### Installation

An executable can be built from source using
and installed locally using:
```bash
go install github.com/ejuju/nubio@latest
```

### Configuration

Your resume is configured using a single JSON file,
usually named `resume.json`.

A `resume.json` file typically contains:
- Contact details
- External links
- Work Experience
- Skills
- Languages
- Education
- Interests
- Hobbies

Check out an example in [/resume.json](/resume.json).

You can check the validity of your `resume.json` using the CLI:
```bash
nubio check-resume-config resume.json
```

### Generating a static website (SSG)

```bash
nubio ssg resume.json static/
```

For an example setup with Caddy on Debian, check out:
[/doc/setup-ssg-caddy-debian.md](/doc/setup-ssg-caddy-debian.md)

### Generating exports via CLI

You can also use Nubio as a CLI to generate static file exports,
for example, to export your config as a PDF:

```bash
nubio export pdf /path/to/resume.json /path/to/output
```

Supported export formats are:
- `html`
- `pdf`
- `json`

### Running as HTTP(S) server

First, you'll need to configure a `server.json` file with the necessary information.

Example `server.json` for a local development server:
```json
{
    "address": ":8080",
    "resume_path": "resume.json"
}
```

Example `server.json` for a HTTPS server (running on port 80 and 443):
```json
{
    "tls_dirpath": "tls",
    "tls_email_addr": "contact@mysite.example",
    "resume_path": "resume.json"
}
```

The field `tls_dirpath` indicates a directory where TLS certificates will be stored.
When using HTTPS, the server uses ports `80` and `443` by default,
there's no need to use the `address` field (it is ignored).

Example `server.json` for a HTTP server behind a reverse proxy:
```json
{
    "address": ":8080",
    "true_ip_header": "X-Forwarded-For",
    "resume_path": "resume.json"
}
```

To start the server, run:
```bash
nubio run server.json
```

NB: You can also simply run `nubio run` which by default will look
for a `server.json` file in the current working directory.

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

- Export your resume to PDF: `nubio.ExportPDF(w, resume)`
- Export your resume to HTML: `nubio.ExportHTML(w, resume)`
- Validate your resume configuration: `resume.Check()`
- And more...

Official package documentation is available here:
[pkg.go.dev/github.com/ejuju/nubio/pkg/nubio](https://pkg.go.dev/github.com/ejuju/nubio/pkg/nubio)
