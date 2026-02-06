# bh - Bithost.io CLI

Official command-line interface for [Bithost.io](https://bithost.io).

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
# Configure authentication
bh auth
```

Get your API key from [dashboard.bithost.io/api_keys](https://dashboard.bithost.io/api_keys).

## Modes

### Interactive TUI

Launch the full-screen terminal UI for visual server management:

```bash
bh manage
```

Features:
- Dashboard with server list and account balance
- Create servers with a step-by-step wizard
- Manage SSH keys
- Delete servers with confirmation
- Copy SSH commands to clipboard
- Auto-refresh server status
- Vim-style navigation (j/k, q to quit)

### Scriptable CLI

Use individual commands for scripting and automation:

```bash
bh servers list
bh servers new --name myserver --provider digital_ocean --region 2 --size 5 --image 10 --keys 1,2
bh servers delete 123 --force
```

## Provider Slugs

Providers are identified by slug strings (e.g., `digital_ocean`) rather than numeric
IDs. Use `bh providers` to list available provider slugs:

```bash
$ bh providers
Slug            Name
digital_ocean   DigitalOcean
...
```

The `--provider` flag on commands like `regions`, `sizes`, `images`, and `servers new`
accepts these slugs.

## CLI Commands

### Server Management

```bash
# List active servers
bh servers list

# List all servers including failed
bh servers list --all

# Create a server
bh servers new --name myserver --provider digital_ocean --region 2 --size 5 --image 10 --keys 1,2

# Create a server with backups enabled
bh servers new --name myserver --provider digital_ocean --region 2 --size 5 --image 10 --keys 1,2 --backups

# Delete a server (with confirmation prompt)
bh servers delete <id>

# Delete a server without confirmation
bh servers delete <id> --force
```

### Resource Discovery

```bash
# List available providers (shows slugs)
bh providers

# List regions for a provider
bh regions --provider digital_ocean

# List sizes/plans for a region
bh sizes --provider digital_ocean --region <region-id>

# List OS images
bh images --provider digital_ocean
bh images --provider digital_ocean --arch arm  # filter by architecture (default: x86)
```

### SSH Key Management

```bash
# List SSH keys
bh ssh-keys list

# Add SSH key (interactive prompts)
bh ssh-keys add

# Add SSH key from file
bh ssh-keys add --label "my-key" --file ~/.ssh/id_rsa.pub

# Add SSH key directly
bh ssh-keys add --label "my-key" --key "ssh-rsa AAAA..."

# Add SSH key for a specific provider (default: digital_ocean)
bh ssh-keys add --label "my-key" --file ~/.ssh/id_rsa.pub --provider digital_ocean
```

### Account

```bash
# Show user information (name, email, balance, server limit)
bh user

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
