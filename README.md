# GitNotify

A powerful, configurable, and low-noise notification tool for GitHub organizations. GitNotify receives webhook notifications from GitHub and provides filtered, real-time updates about issues and pull requests.

## Features

- **Webhook Listener**: Receives and validates GitHub webhook events
- **Event Filtering**: Configurable rules to filter which events trigger notifications
- **Multiple Event Types**: Supports issues and pull request events
- **Repository Filtering**: Monitor specific repositories or entire organizations
- **Action Filtering**: Filter by specific actions (opened, closed, reopened, etc.)
- **Secure**: Validates webhook signatures to ensure requests are from GitHub
- **Health Check**: Built-in health check endpoint for monitoring

## Quick Start

### 1. Clone and Build

```bash
git clone <repository-url>
cd gitnotify
go build -o gitnotify cmd/gitnotify/main.go
```

### 2. Configure

Copy the example configuration and customize it:

```bash
cp config.yml.example config.yml
```

Edit `config.yml` with your settings:

```yaml
organization: "your-organization-name"
port: 8080
webhook_secret: "your-webhook-secret-here"

notifications:
  - event_type: "issues"
    actions: ["opened", "closed", "reopened"]
  
  - event_type: "pull_request"
    actions: ["opened", "closed", "reopened", "synchronize"]
```

### 3. Set Up GitHub Webhook

1. Go to your GitHub organization settings
2. Navigate to Webhooks
3. Add a new webhook with:
   - **Payload URL**: `http://your-server:8080/webhook`
   - **Content type**: `application/json`
   - **Secret**: Use the same secret as in your config.yml
   - **Events**: Select "Let me select individual events" and choose:
     - Issues
     - Pull requests

### 4. Run

```bash
./gitnotify
```

The server will start and log notifications to the console.

## Configuration

### Configuration File Structure

```yaml
# Required: GitHub organization to monitor
organization: "your-org"

# Optional: Port to listen on (default: 8080)
port: 8080

# Required: Webhook secret for validation
webhook_secret: "your-secret"

# Optional: GitHub App settings (for future use)
github_app:
  app_id: 123456
  installation_id: 789012
  private_key_path: "/path/to/private-key.pem"

# Optional: Notification rules
notifications:
  - event_type: "issues"           # Event type to monitor
    actions: ["opened", "closed"]   # Specific actions to notify for
    repos: ["repo1", "repo2"]      # Optional: specific repos to monitor
```

### Notification Rules

If no notification rules are specified, all events will be logged. You can create rules to filter:

- **Event Types**: `issues`, `pull_request`
- **Actions**: `opened`, `closed`, `reopened`, `synchronize`, etc.
- **Repositories**: Specific repository names within the organization

### Examples

**Monitor all issues and PRs:**
```yaml
notifications:
  - event_type: "issues"
  - event_type: "pull_request"
```

**Only new issues and PRs:**
```yaml
notifications:
  - event_type: "issues"
    actions: ["opened"]
  - event_type: "pull_request"
    actions: ["opened"]
```

**Monitor specific repositories:**
```yaml
notifications:
  - event_type: "issues"
    repos: ["important-repo", "security-repo"]
```

## API Endpoints

- `POST /webhook` - GitHub webhook endpoint
- `GET /health` - Health check endpoint

## Development

### Project Structure

```
gitnotify/
├── cmd/
│   └── gitnotify/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── github/
│   │   └── handler.go           # Webhook event handling
│   └── server/
│       └── server.go            # HTTP server setup
├── config.yml.example           # Example configuration
├── go.mod                       # Go module definition
└── README.md                    # This file
```

### Building

```bash
go build -o gitnotify cmd/gitnotify/main.go
```

### Running Tests

```bash
go test ./...
```

## Security

- Webhook signatures are validated using GitHub's `X-Hub-Signature-256` header
- Configuration files containing secrets should not be committed to version control
- The `.gitignore` file excludes `config.yml` and key files by default

## Future Enhancements

- Web UI for configuration management
- Email and Slack notifications
- Advanced filtering (file paths, labels, CODEOWNERS)
- GitHub App authentication
- Database storage for notification history
- AI-powered summaries

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license here] 