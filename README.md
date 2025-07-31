# go-smtp-gateway

A local SMTP gateway that forwards email notifications to Telegram. This project creates a bridge between SMTP-based notifications and Telegram's messaging system, allowing you to receive important alerts directly in your Telegram chat.

## Features

- SMTP server that listens for incoming email notifications
- Message forwarding to Telegram using Telegram Bot API
- Local handling of notification delivery
- Configuration options for Telegram bot and chat settings

## Architecture

The gateway consists of two main components:
1. An SMTP server (main.go) that listens for incoming email messages
2. A Telegram message handler that processes and forwards notifications

## Installation

1. Clone this repository:
```bash
git clone https://github.com/yourusername/go-smtp-gateway.git
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the binary:
```bash
go build -o go-smtp-gateway main.go
```

Or using the ansible script:

1. Run the ansible script in deploy folder
```bash
cd deploy
ansible-playbook -i inventory.ini smtp-gateway.yaml
```

## Configuration

Update the following settings before running:
- Telegram Bot Token: Obtain from BotFather in Telegram
- Telegram Chat ID: Your target chat ID for notifications
- SMTP settings: Configure the SMTP port and listening address

## Usage

Run the gateway with:
```bash
./scripts/run.sh
```

You can also run this project using Docker:

```bash
docker build -t go-smtp-gateway .
docker run -p 2525:2525 -e TELEGRAM_BOT_TOKEN=your_bot_token -e TELEGRAM_CHAT_ID=your_chat_id go-smtp-gateway
```

The Dockerfile is included in this repository for easy deployment.


### Testing

```bash
netcat gateway.home 2525
EHLO gateway.home
AUTH PLAIN 
AGdhdGV3YXkAZ2F0ZXdheQ==
MAIL FROM:<root@nas.local>
RCPT TO:<info@test.local>
DATA
Subject: Test

Hello from the SMTP test.
.
```
where

```bash
echo "AGdhdGV3YXkAZ2F0ZXdheQ==" | base64 -d
gatewaygateway
```

### Features

- Handles SMTP incoming messages
- Formats email content for Telegram delivery
- Supports custom message templates
- Local delivery without external services

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Security Considerations

- Make sure to configure proper firewall rules if exposing the SMTP port
- Keep your Telegram bot token secure
- #TODO: Implement TLS for encrypted communication
