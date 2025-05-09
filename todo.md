# Todo

For v1:
- [ ] Add complete setup examples: SSG with Caddy and HTTPS server with Systemd/Debian.
- [ ] Add `dev` CLI command for running local plain HTTP server for rendering resume file without `server.json` file + hot reload.
- [ ] Make sections (education, hobbies and interests, etc.) optional.
- [ ] Inline custom CSS when exporting as single HTML page.

Nice-to-haves:
- [ ] Add 404 page
- [ ] Add panic recovery page
- [ ] Add meta tags (JSON+LD / OG)
- [ ] Support providing HTTP redirects in config (for URL shortener like capabilities)
- [ ] Support alt domains for HTTP(S) server (server config `"domain_alts": "mysite.fr"` redirecting to `"mysite.example"`)

Ideas:
- [ ] Inline custom CSS file in HTML page head to allow simply opening pages without server for local dev
- [ ] Support i18n
- [ ] Support notifying admin by email on internal server error (panic, etc.)
- [ ] Add more builtin export templates (HTML and PDF)
- [ ] Support serving static files from directory (on `/static/*`) (using file that list file paths, URI and corresponding MIME-type)
- [ ] Support blogging / documentation (with Markdown-like files directory)
- [ ] Support custom HTML template
- [ ] Support analytics reports (page visits / UI events?) sent by email
- [ ] Support hot reload (= (re)generate HTML, PDF, etc. on each request) (useful for local dev)
- [ ] Support contact form with email notification (on dedicated page `/contact`)
- [ ] Support IP blocklist in config / or dedicated file.
- [ ] Add global rate limiting middleware
- [ ] Support logging to file (support file rotation / auto-delete after retention period)
- [ ] Smart page breaks for PDF export (to avoid breaking within content for long resumes)
