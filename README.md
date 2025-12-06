# bh - Bithost.io CLI

Official command-line interface for Bithost.io.

## Installation

```bash
# From source
git clone https://github.com/bithostio/bh.git
cd bh
make build
sudo mv bh /usr/local/bin/

# Or using go install
go install github.com/bithostio/bh@latest
```

## Quick Start

```bash
# 1. Configure authentication
bh auth

# 2. Check balance
bh balance

# 3. Create a server (interactive wizard)
bh servers new --interactive
```

Get your API key from [dashboard.bithost.io/api_keys](https://dashboard.bithost.io/api_keys)

## Commands

### Server Management

```bash
# List all servers
bh servers list

# Create a server (interactive wizard)
bh servers new --interactive

# Create a server (with flags)
bh servers new --name myserver --provider 1 --region 2 --size 5 --image 10 --keys 1,2

# Delete a server
bh servers delete <id>
bh servers delete <id> --force  # skip confirmation
```

### Resource Listing

```bash
# List providers
bh providers

# List regions for a provider
bh regions --provider <provider-id>

# List sizes/plans for a region
bh sizes --provider <provider-id> --region <region-id>

# List OS images
bh images --provider <provider-id>
bh images --provider <provider-id> --arch arm  # filter by architecture
```

### SSH Key Management

```bash
# List SSH keys
bh ssh-keys list

# Add SSH key (interactive)
bh ssh-keys add

# Add SSH key from file
bh ssh-keys add --file ~/.ssh/id_rsa.pub --label "my-key" --provider 1

# Add SSH key directly
bh ssh-keys add --key "ssh-rsa AAAA..." --label "my-key" --provider 1
```

### Account

```bash
# Show balance
bh balance

# Configure API key
bh auth
```

## Configuration

Config is stored in `~/.bh/config.json` with 0600 permissions:

```json
{
  "api_key": "your-api-key",
  "api_base_url": "https://dashboard.bithost.io/api/v1/"
}
```

Environment variables (optional overrides):
- `BH_API_KEY` - API key
- `BH_API_URL` - API base URL

## Troubleshooting

**"config not found"** - Run `bh auth`

**"Authentication failed"** - Get a new API key from [dashboard.bithost.io/api_keys](https://dashboard.bithost.io/api_keys)

**"No SSH keys found"** - Add keys with `bh ssh-keys add` or at [dashboard.bithost.io/keys](https://dashboard.bithost.io/keys)

**Low balance** - Top up at [dashboard.bithost.io/billing](https://dashboard.bithost.io/billing)

## Development

```bash
# Build
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Run modernize
make modernize
```

## License

MIT
