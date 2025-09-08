# Contributing to commitgen

Thank you for your interest in contributing to commitgen! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork and clone the repository**

   ```bash
   git clone https://github.com/joaquinalmora/commitgen.git
   cd commitgen
   ```

2. **Set up development environment**

   ```bash
   make dev-setup
   make build
   ```

3. **Run tests**

   ```bash
   make test
   ```

4. **Test your changes**

   ```bash
   make dev  # Build and test
   make doctor  # System health check
   ```

## 🛠️ Development Workflow

### Building and Testing

```bash
# Quick development cycle
make dev              # Build + test
make test-coverage    # Generate coverage report
make clean            # Clean artifacts

# Manual testing
./bin/commitgen suggest --ai --verbose
./bin/commitgen doctor
```

### Code Quality

- **Go formatting**: Code must be `gofmt` formatted
- **Testing**: Add tests for new functionality
- **Documentation**: Update docs for user-facing changes

### Project Structure

```text
cmd/commitgen/        # Main application entry point
internal/
├── cache/           # Caching system
├── config/          # Configuration management
├── diff/            # Git diff processing
├── doctor/          # System diagnostics
├── hook/            # Git hooks
├── prompt/          # Commit message generation
├── provider/        # AI provider integrations
└── shell/           # Shell integration
```

## 🎯 Areas for Contribution

### High Priority

- **Test coverage expansion** (currently 3/8 packages)
- **Error handling improvements**
- **Performance optimizations**
- **Documentation improvements**

### Medium Priority

- **New AI providers** (following the provider interface)
- **Shell integrations** (bash, fish support)
- **Configuration enhancements**

### Ideas Welcome

- **VS Code extension**
- **Additional shell features**
- **Commit template customization**

## 📋 Pull Request Process

1. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow existing code style
   - Add tests for new functionality
   - Update documentation if needed

3. **Test thoroughly**

   ```bash
   make dev
   make test-coverage
   ./bin/commitgen doctor
   ```

4. **Commit with conventional commits**

   ```bash
   git commit -m "feat: add new feature description"
   # or
   git commit -m "fix: resolve issue with ..."
   ```

5. **Push and create PR**

   ```bash
   git push origin feature/your-feature-name
   ```

## 🧪 Testing Guidelines

### Unit Tests

- Place tests in `*_test.go` files
- Use table-driven tests where appropriate
- Mock external dependencies

### Integration Tests

- Test real workflows in `e2e/` directory
- Test shell integration
- Test AI provider integration

### Manual Testing

```bash
# Test basic functionality
echo "test" > test.txt && git add test.txt
./bin/commitgen suggest --ai --verbose

# Test shell integration
make install
# Then test autosuggestions in a new terminal

# Test system health
./bin/commitgen doctor
```

## 📝 Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment information**

   ```bash
   ./bin/commitgen doctor  # Include this output
   go version
   uname -a
   ```

2. **Steps to reproduce**
3. **Expected vs actual behavior**
4. **Relevant logs or error messages**

## 💡 Feature Requests

- Check existing issues first
- Describe the use case and motivation
- Consider implementation impact

## 🔒 Security

- Report security vulnerabilities privately
- Use GitHub's security advisory feature
- Don't include sensitive information in public issues

## 📄 License

By contributing to commitgen, you agree that your contributions will be licensed under the MIT License.
