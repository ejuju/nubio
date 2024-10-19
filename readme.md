# Nubio: Self-hosted resume for developers

Features:
- Define your resume as JSON
- Export your resume as PDF, JSON, plain text, Markdown, or HTML page
- Serve your resume as a website (or generate static website files).
- Single executable
- Auto HTTPS (using `golang.org/x/crypto/acme/autocert`)

Check out a live demo here: [juliensellier.com](https://juliensellier.com)

## Install

```bash
go install github.com/ejuju/nubio@latest
```

## Usage

Serve over HTTP(S):
- Add your `profile.json` file.
- Add your `config.json` file.
- Then start the server: `nubio run config.json`

Generate a static website (SSG):
- Run: `nubio ssg profile.json static/`
