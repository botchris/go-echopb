# EchoPB

Is a simple gRPC client and server for testing protocol connectivity and performance.

## Installation

### Docker Image

You can run the EchoPB server using the official Docker Image:

```bash
# Start the server on port 4040
docker run --rm botchrishub/echopb:latest /bin/server --listen ":4040"

# Communicate with the server using the client CLI
docker run --rm botchrishub/echopb:latest /bin/client --host ":4040" basic "hello world!"
docker run --rm botchrishub/echopb:latest /bin/client --host ":4040" server-stream "hello world!" --count 1000 --interval 250
```
