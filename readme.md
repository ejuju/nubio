# Nubio: Self-hosted resume for developers

Features:
- Define your resume as JSON
- Export your resume as PDF, JSON, plain text, Markdown, or HTML page
- Serve your resume as a website (or generate static website files).
- Single executable
- Auto HTTPS (using `golang.org/x/crypto/acme/autocert`)

Check out a live demo here: [juliensellier.com](https://juliensellier.com)

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

### HTTP server

Example of config.json for production server with HTTPS:
```json
{
    "domain": "mysite.example",
    "tls_dirpath": "tls",
    "tls_email_addr": "contact@mysite.example",
    "profile": "profile.json",
    "pgp_key": "pgp.asc"
}
```

Example of config.json for production server behind reverse proxy:
```json
{
    "domain": "mysite.example",
    "true_ip_header": "X-Forwarded-For",
    "profile": "profile.json",
    "pgp_key": "pgp.asc"
}
```

Example of config.json for local development server:
```json
{
    "address": ":8080",
    "domain": "localhost",
    "profile": "profile.json",
    "pgp_key": "pgp.asc"
}
```

Make sure the `config.json` and `profile.json` are readable.

To start the server, run:
```bash
nubio run config.json
```

Notes:
- You can also simply run `nubio run` which by default will look 
  for `config.json` in its current working directory.
- The `config.json` indicates where the `profile.json` file can be found.

### Generate a static website (SSG)

```bash
nubio ssg profile.json static/
```

Note: Make sure to fill in the `domain` field in your `profile.json` for SSG.

### Use as a library, for example

- Export profile to PDF: `nubio.ExportPDF(w, profile)`
- Make a HTTP handler that serves the PDF export: `httpHandler := nubio.ExportAndServePDF(profile)`
