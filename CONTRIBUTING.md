# Contributing to go-fantasy-pl

Thanks for considering contributing to go-fantasy-pl! This document will help you get started.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YourUsername/go-fantasy-pl
   cd go-fantasy-pl
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/AbdoAnss/go-fantasy-pl
   ```

## Development Setup

1. Ensure you have Go 1.23 or later installed
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run tests:
   ```bash
   go test ./...
   ```

## Making Changes

1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. Make your changes
3. Run tests and ensure they pass
4. Update documentation if needed
5. Commit your changes with clear messages

## Code Standards

- Follow Go best practices and conventions
- Use meaningful variable and function names
- Add comments for complex logic
- Update tests for new functionality
- Handle errors appropriately

## Pull Request Process

1. Update your fork with the latest upstream changes
2. Push your changes to your fork
3. Create a Pull Request from your fork to our main branch
4. Describe your changes in detail
5. Link any relevant issues

## Testing

- Write unit tests for new functionality
- Ensure all tests pass before submitting PR
- Include integration tests where appropriate
- Test edge cases and error conditions

## Documentation

- Update README.md for significant changes
- Add godoc comments to exported functions
- Include examples for new features
- Update examples in the `/examples` directory

## Questions or Issues?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the codebase
- Documentation improvements

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.
