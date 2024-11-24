![GitHub License](https://img.shields.io/github/license/yosebyte/passport)
[![Go Report Card](https://goreportcard.com/badge/github.com/yosebyte/passport)](https://goreportcard.com/report/github.com/yosebyte/passport)
![GitHub Release](https://img.shields.io/github/v/release/yosebyte/passport)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/yosebyte/passport/docker.yml)
![GitHub last commit](https://img.shields.io/github/last-commit/yosebyte/passport)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/yosebyte/passport/latest)

<div align="center">
  <img src="https://cdn.185610.xyz/assets/passport.png" alt="passport">
</div>

<h4 align="center">"Access pass required to pass through port."</h4>

## Overview

**Passport** is a powerful connection management tool that simplifies network tunneling, port forwarding and more. By seamlessly integrating three distinct running modes within a single binary file, Passport bridges the gap between different network environments, redirecting services and handling connections seamlessly, ensuring reliable network connectivity and ideal network environment. Also with highly integrated authorization handling, Passport empowers you to efficiently manage user permissions and establish uninterrupted data flow, ensuring that sensitive resources remain protected while applications maintain high performance and responsiveness.

## Features

- **Unified Operation**: Passport can function as a server, client, or broker, three roles from a single executable file.

- **Authorization Handling**: By IP address handling, Passport ensures only authorized users gain access to sensitive resources.

- **In-Memory Certificate**: Provides a self-signed HTTPS certificate with a one-year validity, stored entirely in memory.

- **Auto Reconnection**: Providing robust short-term reconnection capabilities, ensuring uninterrupted service.

- **Connection Updates**: In scenarios where connection is interrupted, Passport supports real-time connection updates.

- **Service Forwarding**: Efficiently manage and redirect your connections from one service to entrypoints everywhere.

- **Zero Dependencies**: Fully self-contained, with no external dependencies, ensuring a simple and efficient setup.

## Usage

To run the program, provide a URL specifying the mode and connection addresses. The URL format is as follows:

```
server://linkAddr/targetAddr
client://linkAddr/targetAddr
broker://linkAddr/targetAddr
```

Note that only `server` and  `broker` mode support authorization Handling, which you can just add auth entry after `#`. For example:

```
server://linkAddr/targetAddr#authScheme://authAddr/secretPath
broker://linkAddr/targetAddr#authScheme://authAddr/secretPath
```

- **authScheme**: The option allows you to choose between using HTTP or HTTPS.
- **authAddr**: The server address and port designated for authorization handling.
- **secretPath**: The secret endpoint for processing authorization requests.

### Server Mode

- `linkAddr`: The address for accepting client connections. For example, `:10101`.
- `targetAddr`: The address for listening to external connections. For example, `:10022`.

**Run as Server**

```bash
./passport server://:10101/:10022
```

- This command will listen for client connections on port `10101` , listen and forward data to port `10022`.

**Run as Server with authorization**

```bash
./passport server://:10101/:10022#https://hostname:8443/server
```

- The server handles authorization at `https://hostname:8443/server`, on your visit and your IP logged.
- The server will listen for client connections on port `10101` , listen and forward data to port `10022`.

### Client Mode

- `linkAddr`: The address of the server to connect to. For example, `server_ip:10101`.
- `targetAddr`: The address of the target service to connect to. For example, `127.0.0.1:22`.

**Run as Client**

```bash
./passport client://server_hostname_or_IP:10101/127.0.0.1:22
```

- This command will establish link with `server_hostname_or_IP:10101` , connect and forward data to `127.0.0.1:22`.

### Broker Mode

- `linkAddr`: The address for accepting client connections. For example, `:10101`.
- `targetAddr`: The address of the target service to connect to. For example, `127.0.0.1:22`.

**Run as Broker**

```bash
./passport broker://:10101/127.0.0.1:22
```

- This command will listen for client connections on port `10101` , connect and forward data to `127.0.0.1:22`.

**Run as Broker with authorization**

```bash
./passport broker://:10101/127.0.0.1:22#https://hostname:8443/broker
```

- The server handles authorization at `https://hostname:8443/broker`, on your visit and your IP logged.
- The server will listen for client connections on port `10101` , connect and forward data to `127.0.0.1:22`.

## Container Usage

You can also run **Passport** using container. The image is available at [ghcr.io/yosebyte/passport](https://ghcr.io/yosebyte/passport).

To run the container in server mode with or without authorization:

```bash
docker run --rm ghcr.io/yosebyte/passport server://:10101/:10022#https://hostname:8443/server
```

```bash
docker run --rm ghcr.io/yosebyte/passport server://:10101/:10022
```

To run the container in client mode:

```bash
docker run --rm ghcr.io/yosebyte/passport client://server_hostname_or_IP:10101/127.0.0.1:22
```

To run the container in server mode with or without authorization:

```bash
docker run --rm ghcr.io/yosebyte/passport broker://:10101/127.0.0.1:22#https://hostname:8443/broker
```

```bash
docker run --rm ghcr.io/yosebyte/passport broker://:10101/127.0.0.1:22
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Stars
[![Stargazers over time](https://starchart.cc/yosebyte/passport.svg?variant=adaptive)](https://starchart.cc/yosebyte/passport)
