# LaserGit

[![Build](https://github.com/pinpox/lasergit/actions/workflows/build.yml/badge.svg)](https://github.com/pinpox/lasergit/actions/workflows/build.yml)

<div align="center">
  <img src="logo.png" alt="LaserGit Logo" width="250">
</div>

<div align="center">
  <em>A unified interface to manage pull requests using AGit workflow with Gitea.</em>
</div>

## Overview

LaserGit is a command-line tool that provides an interactive terminal user
interface for managing pull requests in Gitea repositories using the AGit
workflow. It allows you to browse, create, and checkout pull requests without
leaving your terminal.

<div align="center">
  <img src="https://github.com/user-attachments/assets/7ffba505-1189-4dab-b8f4-9b803b67fb08" alt="Screenshot" width="600">
</div>


## Features

- 🔍 Browse pull requests interactively
- ✨ Create new pull requests using AGit workflow
- 🔀 Checkout pull requests locally
- 👁️ View pull request details
- 🔄 Refresh pull request list
- ⌨️ Keyboard-driven interface

## Usage

Navigate to your git repository and run:

```bash
lasergit
```

Alternatively, `--repo <path>` can be used to specify the path to your git repository.

### Interactive Commands

- **↑/↓ arrows**: Navigate through the pull request list
- **Enter**: Checkout the selected pull request
- **c**: Create a new pull request
- **v**: View pull request details
- **r**: Refresh the pull request list
- **q/Esc**: Quit the application

## AGit Workflow

This tool leverages the AGit workflow for creating pull requests. AGit allows
you to create pull requests by pushing to a special reference, eliminating the
need to create branches on the remote repository.

For more information about AGit workflow:
- [Gitea AGit Documentation](https://docs.gitea.com/usage/agit)
- [AGit Flow and Git Repo Guide](https://git-repo.info/en/2020/03/agit-flow-and-git-repo/)

## How It Works

1. **List PRs**: The tool fetches and displays all open pull requests from your
   Gitea repository
2. **Create PR**: Uses AGit to push your current branch with special push
   options to create a new pull request
3. **Checkout PR**: Fetches the pull request as a local branch prefixed with
   `agit-<PR-number>`
