# tensordock-cli

A CLI client for https://tensordock.com

## Installation

Grab a build from [Releases](releases)

## Building

```
go build
```

## Usage

### Configuration


```
$ tensordock-cli config --apiKey <YOUR_API_KEY> --apiToken <YOUR_API_TOKEN> [--serviceUrl <SERVICE_URL>]
```

Note: Go to https://console.tensordock.com/api to get your API key and token

### List Servers

```sh
$ tensordock-cli servers list
```

### Get server info

```sh
$ tensordock-cli servers info --server <serverId>
```

### Start/Stop Server

```sh
$ tensordock-cli servers <start|stop> --server <serverId>
```

### Get billing info

```sh
$ tensordock-cli billing
```