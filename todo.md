# Todo

For v1:
- [ ] Add rate limiting middleware
- [ ] Add 404 page
- [ ] Add panic recovery page
- [ ] Add HTTPS with ACME / autocert
- [ ] Support notifying admin by email on internal server error (panic, etc.)
- [ ] Support simple analytics (number of visits) (sent weekly by email to admin)
- [ ] Dont expose PGP key path here in JSON export

Nice-to-haves:
- [ ] Support contact form with email notification (on dedicated page `/contact`)
- [ ] Add JSON+LD and OG tags for SEO
- [ ] Support blogging / documentation (with Markdown-like files directory)
- [ ] Support serving static files from directory (on `/static/*`) (do we want to use stdlib static file server?)
