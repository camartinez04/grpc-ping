# gRPC Ping Service

This repository contains a simple gRPC-based Ping service implemented in Go. It includes a server that responds to streaming ping requests with pongs, and a client that sends a stream of ping messages.

## Author

Carlos Alvarado Martinez

## Prerequisites

Before you begin, ensure you have the following installed:
- Go (1.22 or higher)
- Protocol Buffer Compiler (protoc)

You can install `protoc` from [Protocol Buffers GitHub release page](https://github.com/protocolbuffers/protobuf/releases) or via package managers like `apt` for Ubuntu or `brew` for macOS.

## Installation

Clone the repository:

```bash
git clone https://github.com/camartinez04/grpc-ping.git
cd grpc-ping
```

## Compiling the Protocol Buffers

Navigate to the directory containing the `grpcping.proto` file and compile it using:

```bash
protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative grpcping.proto
```

This command generates Go code for both the gRPC services and the message types defined in the `grpcping.proto`.

## Download the binaries

### Client
```bash
curl -L -o client-bin-linux.tar.gz https://github.com/camartinez04/grpc-ping/raw/main/bin/client-bin-linux.tar.gz
```

### Server
```bash
curl -L -o server-bin-linux.tar.gz https://github.com/camartinez04/grpc-ping/raw/main/bin/server-bin-linux.tar.gz
```

## Building the Server and Client

### Building the Server

To build the server, navigate to the `server` directory and run:

```bash
cd server
go build -o server
```

This command compiles the server code and outputs an executable named `server`.

### Building the Client

To build the client, navigate to the `client` directory and run:

```bash
cd client
go build -o client
```

This command compiles the client code and outputs an executable named `client`.

## Running the Server

Execute the server binary from within the `server` directory:

```bash
./server
Usage of ./server:
  -port string
      GRPC Server Port
  -server string
      GRPC Server Address
 ```

Optionally, you can specify the address, port and delay seconds to use:

```bash
export GRPC_PING_SERVER="192.168.10.23"
export GRPC_PING_PORT="17002"
./server
```

This will start the server, listening on the default port (e.g., 50051).

## Running the Client

```bash
./client
Usage of ./client:
  -delay string
    	Delay in seconds between pings
  -port string
    	GRPC Server Port
  -server string
    	GRPC Server Address
```
Execute the client binary from within the `client` directory to start sending ping requests, 
you have to specify the address and port of the GRPC Server, and delay seconds to use:

```bash
export GRPC_PING_SERVER="192.168.10.23"
export GRPC_PING_PORT="17002"
export GRPC_PING_DELAY="10"
./client
```

The client will connect to the server, send a series of ping messages, and print the pong responses.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for any improvements or issues you find.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.