# Todo

For v1:
- [ ] Add 404 page
- [ ] Add panic recovery page
- [ ] Add HTTPS with ACME / autocert
- [ ] Support notifying admin by email on internal server error (panic, etc.)
- [ ] Support simple analytics (number of visits) (sent weekly by email to admin)

Nice-to-haves:
- [ ] Support logging to file (and rotate file)
- [ ] Support IP blocklist in config / or dedicated file.
- [ ] Add global rate limiting middleware
- [ ] Support contact form with email notification (on dedicated page `/contact`)
- [ ] Add JSON+LD and OG tags for SEO
- [ ] Support blogging / documentation (with Markdown-like files directory)
- [ ] Support serving static files from directory (on `/static/*`) (do we want to use stdlib static file server?)
- [ ] Support adding PGP key (and serve on `/pgp.asc`, don't provide key path in profile.json, to prevent leak in JSON export)
      OR rely on static file server support to provide PGP key (and add corresponding link in profile?)
