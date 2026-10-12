---
name: xup
description: Upload, download, list, and remove objects from the user's configured S3-compatible bucket using xup and xdown.
---

# xup / xdown

Use this skill when the user wants to manage files in their already configured S3-compatible bucket.

## Setup

Install the **single `xup` binary** (or `xup.exe`) in a writable directory on PATH. The program determines upload (`xup`) versus download (`xdown`) mode from `os.Args[0]`. `xup setup` automatically creates an `xdown` symlink beside `xup` (on Windows it falls back to a hard link when symlinks are unavailable). Do not look for or build a separate xdown binary. If an alias already exists for a different file, setup refuses to overwrite it; if the installation directory is read-only, move xup to a writable directory or arrange the alias manually. If `xup` is not configured, ask the user to run `xup setup` locally or offer headless setup with environment variables. NEVER request that the user paste access keys into chat. Do not print config files, environment variables, or credentials.

Interactive: `xup setup` (password input is hidden on a terminal; config is saved and xdown alias created).

Headless: set `XUP_ENDPOINT`, `XUP_BUCKET`, `XUP_ACCESS_KEY_ID`, `XUP_SECRET_ACCESS_KEY`; optionally `XUP_REGION`, `XUP_SESSION_TOKEN`, `XUP_PATH_STYLE` (`true` by default), then run `xup setup --headless` (also creates the alias). The optional `XUP_CONFIG` sets the configuration file path.

## Commands

- Upload: `xup ./somefile` (stores under `somefile`); `xup ./somefile path/in/bucket` for an explicit object key.
- Download: `xdown somefile` (saves in current directory); `xdown path/in/bucket ./local-name` for an explicit destination.
- List: `xup ls` or `xup ls prefix/`.
- Delete: `xup rm path/in/bucket`.
- Help: `xup --help`.

## Safety

Ask for confirmation before deleting a remote object. Downloads never replace an existing local file. Uploads MAY replace an existing remote object with the same key; clarify or confirm when that would be unsafe. Remote object keys are case-sensitive. Do not mistake the public GitHub repository for a place to store S3 credentials. The CLI configuration is private local JSON, not a file to commit.
