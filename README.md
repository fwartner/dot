# dot

A robust CLI tool to manage, install, and synchronize your dotfiles across multiple systems. The tool supports various Linux distributions (Debian, Ubuntu, Fedora, Arch) and macOS, ensuring your development environment is consistent everywhere.

---

## Features

- **Cross-Platform Compatibility**: Works seamlessly on Linux (Debian, Ubuntu, Fedora, Arch) and macOS.
- **Tool Installation**: Automatically installs required tools like `git`, `zsh`, `neovim`, etc., based on your configuration.
- **Dotfiles Synchronization**:
  - Clones your dotfiles repository.
  - Manages symlinks using [GNU Stow](https://www.gnu.org/software/stow/).
  - Pulls and pushes updates to your dotfiles repository.
- **Customizable**:
  - Configure tools, dotfiles repository, and installation preferences using `config.yml`.
- **Extensible CLI**: Built with [Cobra](https://github.com/spf13/cobra), offering modular commands with flags.

---

## Installation

### Prerequisites

Ensure the following are installed on your system:

- **Git**: For cloning your dotfiles repository.
- **GNU Stow**: For managing symlinks.

### Build from Source

1. **Clone the Repository**

   ```bash
   git clone https://github.com/fwartner/dot.git
   cd dot
   ```

2. **Initialize `go.mod`**

   If `go.mod` does not exist (e.g., for fresh forks), initialize it:

   ```bash
   go mod init github.com/fwartner/dot
   go mod tidy
   ```

   This will create a `go.mod` file and resolve all dependencies.

3. **Build the Binary**

   ```bash
   go build -o dotfiles
   ```

4. **Install Locally**

   Place the binary in your `$PATH` for easier usage:

   ```bash
   mv dotfiles /usr/local/bin/
   ```

Alternatively, download a prebuilt binary from the [Releases](https://github.com/fwartner/dot/releases) page.

---

## Configuration

Create a `config.yml` file in one of the following locations:
- `./config.yml` (current directory)
- `~/.config/dotfiles/config.yml`
- `~/.dotfiles-config.yml`

Example configuration:

```yaml
# Repository URL for your dotfiles
dotfiles_repo: "https://github.com/fwartner/dotfiles-arch.git"

# Directory where the dotfiles will be cloned
dotfiles_dir: "~/dotfiles"

# List of tools to install
tools:
  - git
  - stow
  - zsh
  - neovim
  - tmux
```

---

## Usage

The tool provides several commands to manage your dotfiles. Each command supports additional flags for customization.

### 1. Install Dependencies

```bash
dot install [--skip-zsh]
```

- Installs tools specified in `config.yml`.
- Optionally skips installing Oh My Zsh with the `--skip-zsh` flag.

### 2. Setup Dotfiles

```bash
dot setup [--dir <custom-directory>]
```

- Clones the dotfiles repository.
- Manages symlinks using GNU Stow.
- Use the `--dir` flag to specify a custom directory for dotfiles.

### 3. Pull Updates

```bash
dot pull [--fetch-only]
```

- Pulls the latest changes from your dotfiles repository.
- Use the `--fetch-only` flag to fetch updates without merging them.

### 4. Push Changes

```bash
dot push [--message <commit-message>]
```

- Pushes local changes to the dotfiles repository.
- Use the `--message` flag to specify a custom commit message.

---

## Development

### Prerequisites

- **Go**: Version 1.19 or later.

### Project Structure

```
dot/
├── cmd/            # Command implementations
│   ├── install.go
│   ├── setup.go
│   ├── pull.go
│   ├── push.go
├── utils/          # Reusable utilities
│   └── utils.go
├── main.go         # Entry point
├── config.yml      # Configuration file
```

### Run Locally

```bash
go run main.go <command>
```

---

## Contributing

Contributions are welcome! Feel free to fork the project and submit a pull request.

1. Fork the repository.
2. Create a new feature branch:
   ```bash
   git checkout -b feature/new-feature
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add new feature"
   ```
4. Push to the branch:
   ```bash
   git push origin feature/new-feature
   ```
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgements

- [Cobra CLI Library](https://github.com/spf13/cobra)
- [GNU Stow](https://www.gnu.org/software/stow/)
- [Oh My Zsh](https://ohmyz.sh/)
