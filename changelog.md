# Changelog

## v0.7.0
- Remove "Markdown" and "plain text" exports formats.
- Support resume description (field `description`).
- Improve website style.
- Move skills to top of website.
- Make resume `organization` field optional.

## v0.6.4
- Add default accent color and improve styling

## v0.6.3
- Resume sections "Interests" and "Hobbies" are now optional (not rendered if empty)

## v0.6.0
- Breaking change: `config.json` is now split between `resume.json` and `server.json` files.
- Breaking change: CLI command `check` is now split in `check-server-config` and `check-resume-config`
- Breaking change: Resume website link is included by default in `pdf`, `json`, `txt` and `md` exports
  (not included in `html` export).
- `contact` fields moved to the root of the JSON config.

## v0.5.0
- Breaking change: Move config fields `custom_css`, `custom_css_path`, `pgp_key`, `pgp_key_path` inside `resume` object.
- Breaking change: CLI commands `export` replaces `pdf` and now supports more format types.
- Breaking change: config field `name_slug` renamed to `slug`.
- JSON export is now generated using a text-based template (allowing control over which fields are exposed).

## v0.4.2

- Breaking change: `profile.json` and `server.json` files merged in single `config.json` file.
- Breaking change: CLI commands `check-profile` and `check-server` merged into single command `check`.
- Breaking change: resume is embedded in the `config.json`'s `resume` field.
- Breaking change: resume JSON field `experiences` has been renamed to `work_experience`.
- Custom CSS is now supported.
- "Resume" is now used instead of "profile" for clarity.
- Resume exports are opened in new tab instead of downloaded.
