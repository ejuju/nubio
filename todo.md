# Todo

For v1:
- [ ] Add HTTPS with ACME / autocert
- [ ] CI/CD doc (script for deployment over SSH, systemd daemon file template)

Nice-to-haves:
- [ ] Support logging to file (and rotate file)
- [ ] Support IP blocklist in config / or dedicated file.
- [ ] Support contact form with email notification (on dedicated page `/contact`)
- [ ] Add global rate limiting middleware
- [ ] Add JSON+LD and OG tags for SEO
- [ ] Support blogging / documentation (with Markdown-like files directory)
- [ ] Support serving static files from directory (on `/static/*`)
      Using file that list file paths, URI and corresponding MIME-type
- [ ] Support adding PGP key (and serve on `/pgp.asc`, don't provide key path in profile.json, to prevent leak in JSON export)
      OR rely on static file server support to provide PGP key (and add corresponding link in profile?)
- [ ] Add 404 page
- [ ] Add panic recovery page
- [ ] Support notifying admin by email on internal server error (panic, etc.)
- [ ] Support simple analytics (number of visits) (sent weekly by email to admin)
- [ ] Support providing HTTP redirects in config (for URL shortener like capabilities)
- [ ] Add more builtin HTML and PDF templates
- [ ] Support custom templates (HTML/MD/TXT)
