---
name: xup
description: Upload, download, list, and remove objects from the user's configured S3-compatible bucket using xup and xdown.
---

# xup / xdown

Use this skill when the user wants to manage files in their already configured S3-compatible bucket.

## Setup

Ensure the binaries `xup` and `xdown` are installed and on PATH. If `xup` is not configured, ask the user to run `xup setup` locally or offer headless setup with environment variables. NEVER request that the user paste access keys into chat. Do not print config files, environment variables, or credentials.

Interactive: `xup setup` (password input is hidden on a terminal).

Headless: set `XUP_ENDPOINT`, `XUP_BUCKET`, `XUP_ACCESS_KEY_ID`, `XUP_SECRET_ACCESS_KEY`; optionally `XUP_REGION`, `XUP_SESSION_TOKEN`, `XUP_PATH_STYLE` (`true` by default), then run `xup setup --headless`. The optional `XUP_CONFIG` sets the configuration file path.

## Commands

- Upload: `xup ./somefile` (stores under `somefile`); `xup ./somefile path/in/bucket` for an explicit object key.
- Download: `xdown somefile` (saves in current directory); `xdown path/in/bucket ./local-name` for an explicit destination.
- List: `xup ls` or `xup ls prefix/`.
- Delete: `xup rm path/in/bucket`.
- Help: `xup --help`.

## Safety

Ask for confirmation before deleting a remote object. Downloads never replace an existing local file. Uploads MAY replace an existing remote object with the same key; clarify or confirm when that would be unsafe. Remote object keys are case-sensitive. Do not mistake the public GitHub repository for a place to store S3 credentials. The CLI configuration is private local JSON, not a file to commit.
