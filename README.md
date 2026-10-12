# xup

Minimal Go CLI for an S3-compatible bucket: upload via `xup`, download via `xdown`, list and delete via `xup`.

## Install

Download the build for your OS and CPU under **Actions → Build binaries → latest successful push → Artifacts**. Extract the single `xup` executable (or `xup.exe`) into a **writable directory on your `PATH`**. Run `xup setup` to configure S3 and automatically create `xdown` beside it. Build locally with Go 1.23+:

```sh
go build -o xup ./cmd/xup
```

## Configure

Run `xup setup` for interactive terminal setup. Setup creates a symlink named `xdown` next to the installed `xup` (on Windows it falls back to a hard link if symlinks are unavailable). This is **one binary**, which selects upload or download mode based on the invoked name (`os.Args[0]`). Rerunning setup accepts an alias to the same binary and never overwrites a different existing file. If the directory is read-only, configuration is saved but setup returns an error until you install the executable somewhere writable. Enter endpoint (`https://s3.example.com`), bucket, region (or `-` to omit), access key ID and secret access key. The secret is not echoed on a terminal. Path-style defaults to yes (typical for S3-compatible providers). The bucket must already exist.

In a headless environment, supply environment variables through your shell/secret manager (not command arguments):

```sh
export XUP_ENDPOINT='https://s3.example.com'
export XUP_BUCKET='my-bucket'
export XUP_ACCESS_KEY_ID='...'
export XUP_SECRET_ACCESS_KEY='...'
# optional: XUP_REGION, XUP_SESSION_TOKEN, XUP_PATH_STYLE=false, XUP_CONFIG=/path/to/config.json
xup setup --headless
```

Credentials are saved in your user config directory as `xup/config.json` (e.g. `~/.config/xup/config.json` on Linux), with `0600` file permissions and `0700` directory permissions when newly created. Do not commit or share that file. An `http://` endpoint is supported for trusted environments; HTTPS is recommended. For virtual-host bucket addressing, use `XUP_PATH_STYLE=false` during headless setup.

## Commands

```sh
xup somefile                       # upload as key "somefile"
xup ./local/file.bin backups/a.bin # upload under explicit key
xdown somefile                     # download to ./somefile
xdown backups/a.bin ./a.bin        # download to chosen file
xup ls                             # list all objects (size and key)
xup ls backups/                    # list keys with prefix
xup rm backups/a.bin               # delete exact key
xup --help
```

Upload will overwrite a remote object with the same key. Download refuses to overwrite a local file, and removes a partially downloaded file if the transfer fails. The optional `skills/xup/SKILL.md` teaches compatible AI agents to use these commands and the alias setup safely.
