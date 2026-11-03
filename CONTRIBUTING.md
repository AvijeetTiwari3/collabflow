# Contributing to CollabFlow

Thank you for contributing to CollabFlow! We aim to build a battle-tested, high-concurrency collaborative workspace suite.

## Development Workflow

1. **Fork and Clone**:
   ```bash
   git clone https://github.com/<your-username>/collabflow.git
   cd collabflow
   ```
2. **Launch Stack**:
   ```bash
   make up
   ```
3. **Run Tests**:
   Ensure all Go and frontend unit tests pass with zero race conditions:
   ```bash
   make test
   ```
4. **Submit PR**:
   Open a pull request describing changes, architectural trade-offs, and verification steps.
