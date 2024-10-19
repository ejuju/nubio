# Nubio: Online profile for developers

## Features

- Define your profile as JSON
- Export your profile as PDF, JSON, plain text, Markdown, or HTML page
- Serve your profile as a website (or generate static website files).
- Single executable
- Auto HTTPS (get and renew certs using ACME)

## 3rd-party Go dependencies

- For PDF generation: [github.com/go-pdf/fpdf](https://github.com/go-pdf/fpdf)
- For HTTPS/ACME support: [golang.org/x/crypto](https://golang.org/x/crypto)

## Live demo

- HTML: [juliensellier.com](https://juliensellier.com/)
- PDF [juliensellier.com/profile.pdf](https://juliensellier.com/profile.pdf)
- JSON [juliensellier.com/profile.json](https://juliensellier.com/profile.pdf)
- TXT [juliensellier.com/profile.txt](https://juliensellier.com/profile.pdf)
- Markdown [juliensellier.com/profile.md](https://juliensellier.com/profile.pdf)

## Installation

```bash
go install github.com/ejuju/nubio@latest
```


## Usage

### Create your `profile.json` file

Your `profile.json` groups relevant information for your public profile.
This includes:
- Contact details
- External links
- Work experiences
- Skills
- Languages
- Education
- Interests
- Hobbies

Here's a JSON template you can fill in with your information:
```json
{
    "name": "",
    "contact": {"email_address": ""},
    "links": [
        {"label": "", "url": ""},
    ],
    "experiences": [
        {
            "from": "",
            "to": "",
            "title": "",
            "organization": "",
            "location": "",
            "description": "",
            "skills": [""]
        }
    ],
    "skills": [
        {"title": "", "tools": [""]},
    ],
    "languages": [
        {"label": "", "proficiency": ""},
    ],
    "education": [
        {
            "from": "",
            "to": "",
            "title": "",
            "organization": ""
        }
    ],
    "interests": [""],
    "hobbies": [""]
}
```

Note: See `example.profile.json` for a realistic example.

### Run as HTTP(S) server

First, you'll need to configure your `server.json` file.

Example for a HTTPS server:
```json
{
    "domain": "mysite.example",
    "tls_dirpath": "tls",
    "tls_email_addr": "contact@mysite.example",
    "profile": "profile.json"
}
```

The field `tls_dirpath` indicates a directory where TLS certificates will be stored.
When using HTTPS, the server uses ports `80` and `443` by default,
there's no need to use the `address` field (it is ignored).

Example for a HTTP server behind a reverse proxy:
```json
{
    "domain": "mysite.example",
    "address": ":8080",
    "true_ip_header": "X-Forwarded-For",
    "profile": "profile.json"
}
```

Example for a local development server:
```json
{
    "address": ":8080",
    "domain": "localhost",
    "profile": "profile.json"
}
```

To start the server, run:
```bash
nubio run server.json
```

Notes:
- You can also simply run `nubio run` which by default will look 
  for `server.json` in its current working directory.
- The `server.json` indicates where the `profile.json` file can be found.

### Generate a static website (SSG)

```bash
nubio ssg profile.json static/
```

Note: Make sure to fill in the `domain` field in your `profile.json` for SSG.

### Embed in your Go program

- Export profile to PDF: `nubio.ExportPDF(w, profile)`
- Export profile to HTML: `nubio.ExportHTML(w, profile)`
- Profile type definition: see `nubio.Profile`
- And more...

Official package documentation is available here:
[pkg.go.dev/github.com/ejuju/nubio/pkg/nubio](https://pkg.go.dev/github.com/ejuju/nubio/pkg/nubio)
