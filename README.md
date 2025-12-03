# bh - bithost.io CLI

Official command-line interface for bithost.io.

## Installation

### From Source

```bash
git clone https://github.com/bithostio/bh.git
cd bh
go build -o bh .
sudo mv bh /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/bithostio/bh@latest
```

## Quick Start

### 1. Get Your API Key

Get your API key from [https://dashboard.bithost.io/api_keys](https://dashboard.bithost.io/api_keys)

### 2. Configure Authentication

```bash
bh auth
```

Enter your API key when prompted.

### 3. Check Your Balance

```bash
bh balance
```

### 4. Create a Server

```bash
bh create
```

Follow the interactive wizard to create a new server.

## Commands

### Authentication

```bash
bh auth
```

Configure your bithost.io API key. The key is stored securely in `~/.bh/config.json` with 0600 permissions.

### Balance

```bash
bh balance
```

Display your current account balance. If your balance is below $5, you'll see a warning with a link to top up.

### List Servers

```bash
bh servers
# or
bh ls
# or
bh list
```

Display a table of all your servers with their ID, name, status, IP address, cost, and provider.

### Create Server

```bash
bh create
```

Launch an interactive wizard to create a new server. The wizard will guide you through:

1. **Provider Selection** - Choose your cloud provider
2. **Region Selection** - Select the geographic region
3. **Size/Plan Selection** - Choose server specifications and pricing
4. **Operating System** - Select your OS image
5. **SSH Keys** - Select which SSH keys to add (can select multiple or all)
6. **Backups** - Enable/disable automatic backups (+20% cost)
7. **Server Name** - Choose a name for your server
8. **Confirmation** - Review and confirm your configuration

### Delete Server

```bash
bh delete <server-id>
# or
bh rm <server-id>
```

Delete a server by ID. You'll be prompted for confirmation unless you use the `-f` flag.

**Skip confirmation:**
```bash
bh delete <server-id> -f
```

## Configuration

### Config File

Configuration is stored in `~/.bh/config.json`:

```json
{
  "api_key": "your-api-key",
  "api_base_url": "https://dashboard.bithost.io/api/v1"
}
```

### Environment Variables

You can override the config file with environment variables:

- `BH_API_KEY` - Override configured API key
- `BH_API_URL` - Override API base URL (default: https://dashboard.bithost.io/api/v1)

**Example:**

```bash
export BH_API_KEY="your-api-key"
bh balance
```

## Troubleshooting

### "config not found" error

Run `bh auth` to configure your API key.

### "Authentication failed" error

Your API key may be invalid or expired. API keys are valid for 3 months by default.
Get a new key from [https://dashboard.bithost.io/api_keys](https://dashboard.bithost.io/api_keys) and run `bh auth` again.

### "No SSH keys found" error

You need to add at least one SSH key before creating servers. Add keys at [https://dashboard.bithost.io/keys](https://dashboard.bithost.io/keys).

### Low balance warning

If your balance is below $5, you'll see a warning. Top up at [https://dashboard.bithost.io/billing](https://dashboard.bithost.io/billing).

## Development

### Building

```bash
go build -o bh .
```

### Running Tests

```bash
go test ./...
```

## License

MIT
