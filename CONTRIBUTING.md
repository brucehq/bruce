# Contributing to Bruce

Thank you for considering contributing to the **Bruce.Tools** open-source project! We value your time and effort and want to make the process as smooth as possible. Whether you’re fixing bugs, proposing new features, or improving documentation, your contributions are welcome.

This document outlines the guidelines for contributing to the `brucehq/bruce` repository.

---

## Getting Started

### 1. Reporting Issues
We use GitHub Issues to track bugs, feature requests, and questions. To help us address issues quickly:
- Search [existing issues](https://github.com/brucehq/bruce/issues) to ensure your issue hasn’t already been reported.
- Include clear steps to reproduce the problem, expected behavior, and actual results.
- Attach relevant logs, screenshots, or error messages.

### 2. Submitting Feature Requests
If you have ideas to enhance Bruce.Tools, please:
- Open a **Feature Request** issue and provide detailed use cases.
- Explain the problem your feature would solve and its potential benefits.

---

## Contributing Code

### 1. Prerequisites
Before you begin coding:
1. Familiarize yourself with the Bruce.Tools ecosystem via the [documentation](https://bruce.tools/docs).
2. Install any dependencies and tools required to build and test the code locally (e.g., Docker, Go, etc.).

### 2. Fork & Clone the Repository
1. **Fork** the repository to your GitHub account.
2. Clone the forked repository locally:
```bash
git clone https://github.com/your-username/bruce.git
cd bruce
```

### 3. Branch Strategy
Create a branch for your changes:

Use a descriptive name for your branch, such as fix-issue-123 or feature-improved-logging.
```bash
git checkout -b <branch-name>
```

### 4. Write Clean & Documented Code
Follow the Coding Style: Adhere to Go best practices and Bruce.Tools' conventions.
Include Documentation: Update relevant documentation (e.g., README.md or usage guides) if your changes affect the user experience.
Write Tests: Ensure your contributions include unit tests, integration tests, or end-to-end tests when applicable.

### 5. Commit Your Changes
Write clear and concise commit messages:

```bash
git commit -m "Fix issue #123: Improve error handling in X"
```
Push your branch to your forked repository:

```bash
git push origin <branch-name>
```
### 6. Open a Pull Request
Submit a pull request (PR) to the main branch of the official repository:

Navigate to your fork on GitHub.
Click Compare & Pull Request.
Provide a detailed description of your changes, including the problem being solved and any relevant context.
Reference any related issues (e.g., Fixes #123).

---

## Pull Request Guidelines
### 1. Review Process
* Your PR will be reviewed by project maintainers. Please respond promptly to any feedback or requests for changes.
* Ensure your PR passes all checks (e.g., CI/CD pipelines, linting, and tests).
### 2. Code Quality
* Keep changes focused on a single feature, bug fix, or improvement.
* Avoid unrelated modifications.
* Ensure backward compatibility unless introducing a major version update.
### 3. Checklist
Before submitting your PR:
* Code compiles without errors.
* Tests pass locally.
* Code adheres to the project’s coding style.
* Documentation has been updated if necessary.
* * For the documentation site please see this repository: [Bruce Documentation Repository](https://github.com/brucehq/bruce-docs)

---

## Contributor Community
Join our growing community:

Engage with other contributors via [Discussions](https://github.com/brucehq/bruce/discussions).
Follow updates on the use cases & guides section of the website: [Use Cases & Guides](https://bruce.tools/guides).
For private inquiries, contact us at [support@bruce.tools](mailto:support@bruce.tools).

---
## Contributor License Agreement (CLA)
To protect the project and contributors, you may be required to sign a Contributor License Agreement (CLA). By submitting code, you agree to license your contribution under the same terms as the project.

---
## Code of Conduct
We expect all contributors to adhere to our Code of Conduct. Be respectful and professional in all interactions.

## Thank You!
Your contributions help make Bruce better for everyone. We appreciate your effort and look forward to building a thriving community together.

Happy coding! 🚀
