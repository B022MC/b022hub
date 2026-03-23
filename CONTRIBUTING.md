# Contributing

Thanks for contributing to b022hub.

## Before You Start

- Check existing issues and pull requests before starting overlapping work.
- For larger changes, open an issue first so the scope and direction are clear.
- Keep pull requests focused. Separate refactors from behavior changes when possible.

## Development Setup

See [DEV_GUIDE.md](DEV_GUIDE.md) for the current local setup, commands, and contributor workflow.

## Pull Request Checklist

- Run the relevant backend and frontend tests locally.
- Update docs or examples when configuration, API behavior, or deployment steps change.
- Include generated files when the source change requires them, such as Ent output.
- Do not commit secrets, real API keys, populated `.env` files, or private production endpoints.

## Coding Expectations

- Follow the established project structure and naming where possible.
- Prefer small, reviewable commits and clear commit messages.
- Add or update tests when fixing bugs or changing behavior.
- Preserve upstream attribution and keep the repository clearly distinguishable from the upstream Sub2API project.

## Reporting Problems

- Use issues for bug reports, questions, and feature discussions.
- Use private security reporting for vulnerabilities when available. See [SECURITY.md](SECURITY.md).
