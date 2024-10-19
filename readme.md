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

Run server:
- Add your `profile.json` file.
- Add your `config.json` file.
- Then start the server: `nubio run`

Generate a static website (SSG): 
- Create output directory: `mkdir static`
- Then generate the files: `nubio ssg profile.json static`
