> [!IMPORTANT]
> Replaced with https://github.com/sethrylan/slack-reader

# slack-threads

A CLI tool that lists threaded conversations from a Slack channel. Uses cookie-based authentication via [rneatherway/slack](https://github.com/rneatherway/slack) to access the Slack API without requiring a bot token.

## Install

```sh
go install github.com/sethrylan/slack-threads@latest
```

## Usage

```sh
slack-threads -channel <CHANNEL_ID> -domain <TEAM_DOMAIN> [-lookback <DAYS>]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-channel` | Slack channel ID | (required) |
| `-domain` | Slack workspace domain (the `<domain>` in `<domain>.slack.com`) | (required) |
| `-lookback` | Number of days to look back for messages | `7` |

### Example

```sh
slack-threads -channel C01ABCDEF -domain myworkspace -lookback 14
```

Output is a list of Slack URLs linking to each thread parent message found in the channel.

## Authentication

This tool uses cookie-based authentication. Run `slack-threads` and it will prompt for your Slack session cookies. See [rneatherway/slack](https://github.com/rneatherway/slack) for details.

## License

[MIT](LICENSE)
