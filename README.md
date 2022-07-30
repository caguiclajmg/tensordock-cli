# tensordock-cli

A CLI client for [TensorDock](https://tensordock.com)

## Installation

Install using `go install github.com/caguiclajmg/tensordock-cli` or grab a build from [Releases](https://github.com/caguiclajmg/tensordock-cli/releases)

## Build

```
go build
```

## Usage

Add `--help` to any command to get contextual help

### Configuration

```
tensordock-cli config --apiKey api_key --apiToken api_token [--serviceUrl service_url]
```

Go to https://console.tensordock.com/api to get your API key and token

Credentials may also be specified inline with every command using the `--apiKey` and `--apiToken` flags

### List servers

```sh
tensordock-cli servers list
```

### Get server info

```sh
tensordock-cli servers info server_id
```

### Start/Stop/Restart server

```sh
tensordock-cli servers start|stop|restart server_id
```

### Delete server

```sh
tensordock-cli servers delete server_id
```

### Open management dashboard in browser

```sh
tensordock-cli servers manage server_id
```

### Deploy a server

```sh
tensordock-cli servers deploy \
    [--gpuModel gpu_model \]
    [--location location \]
    [--instanceType instance_type \]
    [--gpuCount gpu_count \]
    [--vcpus vcpus \]
    [--storage storage \]
    [--storageClass storage_class \]
    [--ram ram \]
    [--os os \]
    name \
    admin_user \
    admin_pass
```

#### Deploy a GPU Server

```sh
tensordock-cli servers deploy server_name admin_user admin_pass --gpuCount 2 --gpuModel A4000
```

#### Deploy a CPU-only Server

```sh
tensordock-cli servers deploy server_name admin_user admin_pass --instanceType cpu --cpuModel Intel_Xeon_V4
```

#### Convert a server to a CPU instance

```sh
tensordock-cli servers modify server_id --instanceType cpu --cpuModel Intel_Xeon_V4
```

#### Convert a server to a GPU instance

```sh
tensordock-cli servers modify server_id --instanceType gpu --gpuModel Quadro_4000 --gpuCount 2
```

### Get billing info

```sh
tensordock-cli billing
```

### Get GPU stock

```sh
tensordock-cli stock list
```

### Get CPU stock

```sh
tensordock-cli stock list --type cpu
```