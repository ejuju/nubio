# Todo

For v1:
- [ ] Add complete setup examples: SSG with Caddy and HTTPS server with Systemd/Debian.
- [ ] Serve version number on (`/version` / `$ nubio version`)
- [ ] Support alt domains for HTTP(S) server (server config `"domain_alts": "mysite.fr"` redirecting to `"mysite.example"`)

Nice-to-haves:
- [ ] Smart page breaks for PDF export (to avoid breaking within content for long profiles)
- [ ] Support logging to file (support file rotation / auto-delete after retention period)
- [ ] Support IP blocklist in config / or dedicated file.
- [ ] Support contact form with email notification (on dedicated page `/contact`)
- [ ] Add global rate limiting middleware
- [ ] Add meta tags (JSON+LD / OG)
- [ ] Support blogging / documentation (with Markdown-like files directory)
- [ ] Support serving static files from directory (on `/static/*`) (using file that list file paths, URI and corresponding MIME-type)
- [ ] Support i18n
- [ ] Add 404 page
- [ ] Add panic recovery page
- [ ] Support notifying admin by email on internal server error (panic, etc.)
- [ ] Support simple analytics (number of visits) (sent weekly by email to admin)
- [ ] Support providing HTTP redirects in config (for URL shortener like capabilities)
- [ ] Add more builtin HTML and PDF templates
- [ ] Support custom templates (HTML/MD/TXT)
- [ ] Support hot reload (= (re)generate HTML, PDF, etc. on each request) (useful for local dev)
