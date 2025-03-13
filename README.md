# Bappa Game Template Generator

A simple tool for generating new Bappa game projects from templates.

## Overview

Bappa Game Template Generator helps you quickly bootstrap new game projects using the Bappa game engine. It creates a complete project structure with all necessary files and dependencies based on predefined templates.

## Installation

### Option 1: Install with Go (Recommended)

```bash
go install github.com/TheBitDrifter/bappacreate@latest
```

### Option 2: Install from source

```bash
# Clone the repository
git clone https://github.com/TheBitDrifter/bappacreate.git
cd bappacreate

# Build and install the binary
go install
```

## Usage

```bash
bappacreate username/project-name [--template <template-name>]
```

The `username/` prefix is important as it will be used to create the proper Go module path (`github.com/username/project-name`).

### Examples

Create a top-down game (default template):

```bash
bappacreate johndoe/my-awesome-game
```

Create a platformer game:

```bash
bappacreate johndoe/my-platformer --template platformer
```

Create a sandbox game:

```bash
bappacreate johndoe/my-sandbox-world --template sandbox
```

## Available Templates

| Template | Description |
|----------|-------------|
| `topdown` | A top-down perspective game (default) |
| `topdown-split` | A top-down game with split-screen co-op support |
| `platformer` | A simple platformer game |
| `platformer-split` | A platformer game with split-screen co-op support |
| `sandbox` | An open sandbox game environment |

### Template Structure

- **Standard Templates** (`topdown`, `platformer`, `sandbox`): Single-player games with a standard view.
- **Split Templates** (`topdown-split`, `platformer-split`): Games with split-screen co-op multiplayer support, allowing two or more players to play simultaneously on the same screen.

## Project Structure

When you create a new project, the following structure will be generated:

```
my-awesome-game/
├── assets/
│   ├── images/
│   └── sounds/
├── players/        # Only for split-screen co-op templates
├── splitscreen/    # Only for split-screen co-op templates
├── go.mod
├── go.sum
└── main.go
```

## After Creating a Project

After creating your project, navigate to the project directory and run it:

```bash
cd my-awesome-game
go run .
```

## Dependencies

The generator automatically sets up the following dependencies:

- github.com/TheBitDrifter/coldbrew
- github.com/TheBitDrifter/blueprint
- github.com/TheBitDrifter/warehouse

## Contributing

We welcome contributions to add new templates or improve existing ones! Simply add your template to the `templates/` directory and submit a pull request.

## License

[MIT License](LICENSE)
