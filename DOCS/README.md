# CommitGen Documentation

Complete documentation for commitgen - the AI-powered git commit message generator.

## 📚 Documentation Index

### Getting Started

- **[Installation & Setup](SETUP.md)** - Complete installation, configuration, and shell integration guide

### Development

- **[Contributing Guidelines](CONTRIBUTING.md)** - How to contribute to the project
- **[Technical Reference](TECHNICAL.md)** - Architecture, components, and development guide  
- **[Deployment Roadmap](TODO_DEPLOYMENT.md)** - Development progress and deployment checklist
- **[Changelog](CHANGELOG.md)** - Version history and release notes

## 🚀 Quick Start

```bash
# Install via Homebrew (recommended)
brew tap joaquinalmora/tap
brew install commitgen

# Or via Go
go install github.com/joaquinalmora/commitgen/cmd/commitgen@latest

# Interactive setup
commitgen init

# Basic usage
git add .
commitgen suggest
```

## 📖 Key Documents

- **[SETUP.md](SETUP.md)** - Complete installation and configuration guide
- **[TECHNICAL.md](TECHNICAL.md)** - Architecture and development details
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines and development workflow

## 🔧 Configuration

CommitGen supports multiple configuration methods:

1. **Interactive setup**: `commitgen init`
2. **YAML configuration**: `commitgen.yaml` file
3. **Environment variables**: `.env` file
4. **Command line flags**: `--ai`, `--verbose`, etc.

See [SETUP.md](SETUP.md) for detailed configuration options.

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Development setup
- Code style guidelines
- Testing requirements
- Pull request process

## 📋 Project Status

See [TODO_DEPLOYMENT.md](TODO_DEPLOYMENT.md) for current development status and deployment progress.
