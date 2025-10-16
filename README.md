# EchoPB

Is a simple gRPC client and server for testing protocol connectivity and performance.

## Installation

### Docker Image

You can run the server and client using the official Docker Image:

```sh
# Start the server on port 4040
docker run --rm -p 4040:4040 botchrishub/echopb:latest /usr/local/bin/server --listen ":4040"

# Communicate with the server using the client CLI
docker run --rm botchrishub/echopb:latest /usr/local/bin/client --host ":4040" basic "hello world!"
docker run --rm botchrishub/echopb:latest /usr/local/bin/client --host ":4040" server-stream "hello world!" --count 1000 --interval 250
```

### Go Installation

```sh
go install github.com/botchrishub/echopb/cmd/echopb-client@latest
go install github.com/botchrishub/echopb/cmd/echopb-server@latest
```

Then you can run the server and client:

```sh
echopb-server --listen ":4040"
echopb-client --host ":4040" basic "hello world!"
```
