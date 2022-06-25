# tensordock-cli

A CLI client for https://tensordock.com

## Installation

Grab a build from [Releases](https://github.com/caguiclajmg/tensordock-cli/releases)

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

Credentials may also be specified inline with every command using the `--apiKey` and `--apiToken` flags

### List servers

```sh
$ tensordock-cli servers list
```

### Get server info

```sh
$ tensordock-cli servers info --server <serverId>
```

### Start/Stop server

```sh
$ tensordock-cli servers <start|stop> --server <serverId>
```

### Delete server

```sh
$ tensordock-cli servers delete --server <serverId>
```

### Deploy server

```sh
$ tensordock-cli servers deploy \
    [--instanceType <INSTANCE_TYPE> \]
    [--gpuCount <GPU_COUNT> \]
    [--vcpus <VCPUS> \]
    [--storage <STORAGE> \]
    [--storageClass <STORAGE_CLASS> \]
    [--ram <RAM> \]
    [--os <OS> \]
    --adminUser <ADMIN_USER> \
    --adminPass <ADMIN_PASS> \
    --gpuModel <GPU_MODEL> \
    --location <LOCATION> \
    --name <NAME>
```

### Get billing info

```sh
$ tensordock-cli billing
```