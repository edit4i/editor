# Edit4I 🚀 (ed - di - tai)

An AI-powered agentic editor that enhances your coding experience with intelligent context understanding and seamless workflow integration.
An editor that increases your productivity by 4x (4x increase, get it? 😉).

![screenshot-2024-12-03-10-38-50](https://github.com/user-attachments/assets/abc2db3f-5f17-479e-82fc-8618eddd6473)

> This is VERY MUCH in the early stages. Not suitable for production at all. But you can star ⭐ this repo and get notified when we post our first public release!
> Or you can build it yourself, though it's not ready yet.


## Vision 🎯

I'm very frustaded with some editors feature gating a bunch of stuff in their own models, using your own APIs is pratically unusable. I know, I know. They also need to pay the bills, that's fair. But even paying having some errors in my screen and not being able to put my own API key...

We'll, just complain won't solve anything. Let's build something better!

## Stack 🛠️

- **Frontend**
  - ⚡ SvelteJS - Reactive UI framework
  - 📝 Monaco - VS Code's editor component
  - 🌳 Tree-sitter - Robust parsing engine
  
- **Backend**
  - 🚀 Go - Performance-focused backend
  - 🔄 Wails - Native desktop integration (multiplataform)

EditAI combines these technologies to create a high-performance editor with native capabilities while maintaining the flexibility of web technologies. The Go backend enables efficient AI pipeline processing and local workload handling, while SvelteJS and Monaco provide a modern, responsive editing experience.

## Features ✨

- **Automatic Context Generation** 🧠: Intelligently understands your codebase and provides relevant context for smarter coding assistance
- **Treesitter Based Document Symbols** 🌳: Advanced code parsing and symbol recognition for improved code navigation and understanding
- **Git Based Workflow** 🔄: Seamless integration with Git for efficient version control and collaboration
- **Terminal Execution and Debugging** 💻: Built-in terminal support for running commands and debugging directly within the editor
- **Pinned Prompts** 📌: AI convention-based system for consistent and customizable AI interactions

## Roadmap 🛣️

MVP 1:
- [ ] Monaco Editor with Vim mode & Git integration (In progress 🏃)
- [ ] File explorer with search
- [ ] AI Integration (Copilot, OpenAI, Anthropic)
- [ ] Basic context management
- [ ] Integrated terminal with command execution
- [ ] CLI interface for core operations

MVP 2:
- [ ] Multi-file context understanding
- [ ] Dependency graph analysis
- [ ] Web page & documentation parsing
- [ ] Command suggestion & interpretation
- [ ] Task management system
- [ ] External tool integration
- [ ] Custom commands & webhooks

MVP 3:
- [ ] Autonomous AI agents
- [ ] Code semantic search
- [ ] YAML-based tool creation system
- [ ] Tool marketplace & version management
- [ ] Organization knowledge base
- [ ] Team collaboration features
- [ ] CI/CD workflow integration
- [ ] Project configuration (.editai)

You can check the full roadmap [here](./ROADMAP.md).

## Development

### Prerequisites

- Go 1.22+
- Node.js 22+
- Docker (optional, for containerized builds)
- Wails v2.9.2+
- UPX (optional, for production builds)
- NSIS (optional, for Windows builds)

### Building from Source

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/EditAI.git
   cd EditAI
   ```

2. **Local Build**
   ```bash
   # Install dependencies
   make deps

   # Development mode with hot reload
   make dev

   # Development build
   make build-dev

   # Production build
   make build
   ```

3. **Docker Build** (⚠️ Coming Soon)
   > Note: Docker builds are currently being reworked and are temporarily broken. Please use local builds for now.
   ```bash
   # Docker build (currently unavailable)
   docker-compose up
   ```

### Editor Configuration

Create a `.edit4i/config.yaml` file in your home directory to configure the editor.

```yaml
editor:
  theme: "vs-dark"  # or "vs-light", "hc-black"
  fontSize: 14
  tabSize: 2
  wordWrap: true
  lineNumbers: true
  minimap: true
  vim:
    enabled: true
    defaultMode: "normal"  # "normal", "insert", "visual"

```

### Build Targets

- Linux: `make build-linux`
- macOS: `make build-darwin`
- Windows: `make build-windows`

For development and debugging, use `make build-dev` which enables hot reload.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Screenshots

**Source control**


![screenshot-2024-12-03-10-42-12](https://github.com/user-attachments/assets/75de0fac-547a-45d7-bdd1-e31843cefe26)

**Command pallete**

![screenshot-2024-12-03-10-40-10](https://github.com/user-attachments/assets/2d36dd1f-c0a9-4cd4-83a6-f5064c49a735)
> With `Alt + J`, you can use "J" and "K" to move up and down. Yes, I'm a fellow vimmer 😎


**Fuzzy finder**

![screenshot-2024-12-03-11-27-56](https://github.com/user-attachments/assets/e18e7d38-aaf5-43e4-b5a9-6c1e932059eb)


## License

I'm not sure yet. Some of the collaboration features will be probabily licensed under a commercial one when you pass a given amout of users. But we'll always offer a self-hosted version and you'll be always welcome to use your own model (aka, no more "resource_exhausted" errors).
[License details to be added]
