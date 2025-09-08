# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Cache clear functionality (`commitgen cache --clear`)
- Automatic .env file loading for seamless configuration  
- Simplified OpenAI-only integration (removed Ollama complexity)
- Enhanced configuration with smart defaults
- Comprehensive development documentation
- Build automation with Makefile
- MIT License for open-source compliance

### Changed

- Streamlined AI provider selection (OpenAI only)
- Improved environment variable handling
- Simplified configuration variable names
- Enhanced .env.example with better documentation

### Fixed

- Duplicate commit message suggestions (disabled shell integration conflicts)
- Environment variable loading issues
- Configuration file detection
- API key validation workflow
- Cross-platform CI/CD issues (Windows PowerShell compatibility)

### Removed

- Ollama provider support (simplified to OpenAI only)
- Complex provider switching logic
- Redundant configuration options

## [0.1.0] - Initial Development

### Added in 0.1.0

- Core AI-powered commit message generation
- OpenAI integration with gpt-4o-mini support
- High-performance caching system with git hooks
- Shell integration for zsh autosuggestions
- Git workflow integration (prepare-commit-msg hook)
- Heuristic fallback for offline usage
- Doctor command for system diagnostics
- Professional commit message conventions
- Shell installation/uninstallation commands
- Basic CI/CD pipeline
- Documentation structure

### Features

- Sub-100ms commit message retrieval via cache
- Background pre-generation on `git add`
- Ghost text suggestions in terminal
- Smart fallback chain: AI → Cache → Heuristics
- One-command setup experience
- Zero-config defaults with customization options
