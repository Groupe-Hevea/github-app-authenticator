# github-app-authenticator

A simple CLI script that provides you github app installation access tokens.

We use it internally at @Groupe-Hevea to handle some private apps that require interactions with Github APIs.

## Workflow

See https://docs.github.com/en/developers/apps/authenticating-with-github-apps for more details

## Installation

1. Download the latest available release from the release page
2. `chmod +x github-app-authenticator`
3. move the binary in a folder which is in your PATH (like `/usr/local/bin`)

## Commands

* `github-app-authenticator -v` -> get the app version
* `github-app-authenticator {app_id} {private_key_pem_path} {app_installation_id}` -> outputs an installation access token 

## Contributing

### Requirements

* Go 1.15

### Building a snapshot

* `make build`

### Building a release

* `VERSION=1.2.3 make build`
