# EchoPB

Is a simple gRPC client and server for testing protocol connectivity and performance.

## Installation

### Docker Image

You can run the EchoPB server using the official Docker Image:

```bash
# Start the EchoPB server on port 4040
docker run --rm botchrishub/echopb:latest --listen ":4040"

# Run the client against the EchoPB server
docker run --rm botchrishub/echopb:latest /bin/client --host ":4040" basic --message "Hello, EchoPB!"
```