# Cloud CLI Tool

A command-line interface tool for managing cloud resources across different providers (AWS, GCP, Azure).

## Features

- **Multi-Cloud Support**: Unified interface for AWS, GCP, and Azure
- **Resource Management**:
  - Virtual Machines (create, list, start, stop, delete)
  - Storage (buckets and objects management)
- **Output Formats**: Support for table, JSON, and YAML output
- **Configuration Management**: YAML-based configuration with environment variable support

## Installation

```bash
go install github.com/yourusername/cloud-cli@latest
```

## Configuration

Create a configuration file at `~/.cloud-cli.yaml`:

```yaml
default_provider: aws
default_region: us-west-2
providers:
  aws:
    region: us-west-2
    credentials:
      access_key_id: YOUR_ACCESS_KEY
      secret_access_key: YOUR_SECRET_KEY
  gcp:
    region: us-central1
    credentials:
      credentials_file: /path/to/credentials.json
  azure:
    region: eastus
    credentials:
      subscription_id: YOUR_SUBSCRIPTION_ID
      tenant_id: YOUR_TENANT_ID
      client_id: YOUR_CLIENT_ID
      client_secret: YOUR_CLIENT_SECRET
```

## Usage

### Global Flags
- `--config`: Config file (default is $HOME/.cloud-cli.yaml)
- `--provider, -p`: Cloud provider (aws, gcp, azure)
- `--region, -r`: Cloud region
- `--output, -o`: Output format (table, json, yaml)

### Virtual Machines

List VMs:
```bash
cloud vm list
```

Create VM:
```bash
cloud vm create --name myvm --type t2.micro --image ami-12345678
```

Start VM:
```bash
cloud vm start vm-12345678
```

Stop VM:
```bash
cloud vm stop vm-12345678
```

Delete VM:
```bash
cloud vm delete vm-12345678
```

### Storage

List buckets:
```bash
cloud storage list-buckets
```

Create bucket:
```bash
cloud storage create-bucket --name mybucket
```

List objects:
```bash
cloud storage list-objects mybucket
```

Upload object:
```bash
cloud storage upload myobject --bucket mybucket --file /path/to/file
```

Download object:
```bash
cloud storage download myobject --bucket mybucket --output /path/to/save
```

Delete object:
```bash
cloud storage delete-object mybucket myobject
```

## Development

### Project Structure

```
.
├── cmd/
│   ├── root.go
│   ├── vm.go
│   └── storage.go
├── internal/
│   ├── cloud/
│   │   ├── client.go
│   │   ├── aws.go
│   │   ├── gcp.go
│   │   └── azure.go
│   └── config/
│       └── config.go
├── pkg/
│   └── output/
│       └── output.go
├── main.go
└── go.mod
```

### Building from Source

```bash
git clone https://github.com/yourusername/cloud-cli.git
cd cloud-cli
go build
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
