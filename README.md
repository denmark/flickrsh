# flickrsh

A command-line tool for managing your Flickr photos: list albums, upload files, and download photos/videos.

## Installation

Requires Go 1.25+.

```sh
make build
```

This produces a `flickrsh` binary in `bin/`. Alternatively, install it to your `$GOPATH/bin`:

```sh
make install
```

## Setup

Before using `flickrsh`, you need a Flickr API key and secret. Obtain these from the [Flickr App Garden](https://www.flickr.com/services/apps/create/).

Initialize the connection by running:

```sh
flickrsh init --key <API Key> --secret <API Secret>
```

This will:

1. Print an authorization URL. Open it in a browser, sign in to Flickr, and approve access.
2. Prompt you to enter the OAuth confirmation code shown by Flickr.
3. Save your credentials to `~/.flickrsh.yml` (`flickrsh.yml` in your home directory on Windows).

If a configuration file already exists, `init` will ask for confirmation before overwriting it.

## Usage

### List albums

```sh
flickrsh album
```

Prints the ID and title of every album (photoset) in your Flickr account.

### Upload photos/videos

```sh
flickrsh upload [flags] <file> [file...]
```

Flags:

| Flag | Description |
| --- | --- |
| `--tag <tag>` | Tag to apply to uploaded files. Repeat for multiple tags. |
| `--album <name>` | Album to add the uploaded files to. |
| `--public` | Make uploaded files visible to the public. |
| `--family` | Make uploaded files visible to family. |
| `--friend` | Make uploaded files visible to friends. |
| `--searchable` | Make uploaded files available to public search. |

Example:

```sh
flickrsh upload --tag vacation --tag beach --public photo1.jpg photo2.jpg
```

Failed uploads are automatically retried up to 3 times before being skipped.

### Download photos/videos

```sh
flickrsh download [flags]
```

Flags:

| Flag | Description |
| --- | --- |
| `--album <id>` | ID of the album to download photos from (see `flickrsh album` for IDs). |
| `--minheight <n>` | Minimum image height to download; picks the smallest available size that meets this threshold. |
| `--dir <path>` | Directory to save downloaded files to (default: current directory). |

Example:

```sh
flickrsh download --album 72177720123456789 --minheight 1080 --dir ./downloads
```

Files already present in the target directory are skipped.

## Configuration

`flickrsh` stores its Flickr credentials in `~/.flickrsh.yml`:

```yaml
user_id: ...
api_key: ...
api_secret: ...
oauth_token: ...
oauth_token_secret: ...
```

You generally shouldn't need to edit this file by hand; use `flickrsh init` to (re)create it.

## Development

```sh
make build   # build the binary
make test    # run tests
make vet     # run go vet
make fmt     # format code
make tidy    # tidy go.mod/go.sum
make clean   # remove build output
```

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
