# spotcapture
Capture network packets until a stop event is received and upload PCAP file with captured packages to an S3 bucket.

## Features
* Continuous packet capture on specified network interface
* Ring buffer for efficient memory management
* PCAP file generation from captured packets
* Automatic upload to AWS S3
* HTTP proxy support
* Multiple stop event triggers

## Stop Events
* SIGTERM signal
* TCP/UDP connection on specified port

## Configuration
Create a `config.yml` file in the project root:

```yaml
interface: eth0          # Network interface to capture on
port: 5000               # Port to listen for stop events
key: "<aws_access_key>"  # AWS Access Key ID
secret: "<aws_secret>"   # AWS Secret Access Key
region: "eu-central-1"   # AWS region
bucket: "s3-bucket-name" # S3 bucket name
proxy: "http://proxy"    # Optional HTTP proxy URL
```

## Build

```bash
CGO_ENABLED=0 go build -o spotcapture cmd/main.go
```

## How It Works
1. Loads configuration from `config.yml`
2. Starts packet capture on the specified network interface
3. Stores packets in a ring buffer (max 1000 packets)
4. Listens for TCP/UDP connections on the specified port
5. On stop event, converts captured packets to PCAP format
6. Uploads PCAP file to S3 bucket with timestamp

## Requirements
* Go 1.16 or later
* Network interface access (may require elevated privileges)
* AWS credentials with S3 write permissions
