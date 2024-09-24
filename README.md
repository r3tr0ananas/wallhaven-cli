<div align="center">

# 🖼️ wallhaven-cli

<sub>Search and download wallpapers from [wallhaven](https://wallhaven.cc).</sub>

</div>

## 📦 Installation

### Prerequisites

Before installing `wallhaven-cli`, make sure your system meets the following requirements:

- **Supported platform**:  
  - Linux
- **[Go](https://go.dev/)**:  
  Required to build this project from the source.
- **[fzf](https://github.com/junegunn/fzf?tab=readme-ov-file#installation)**:  
  Used for the selection menu. This is **required**.
- **[kitty](https://github.com/kovidgoyal/kitty)**:  
  Currently, `kitty` is **required** as the terminal emulator.

### Installing

Once all prerequisites are met, you can install `wallhaven-cli` using one of the following methods:

#### 1. Install via `go install`

```bash
go install github.com/r3tr0ananas/wallhaven-cli/wallhaven@latest
```

#### 2. AUR

You can install `wallhaven-cli` via the AUR:
```bash
yay -S wallhaven-cli
```

If you use a different AUR helper then replace `yay`.

#### 3. Download the Latest Release

Alternatively, you can download the [latest release](https://github.com/r3tr0ananas/wallhaven-cli/releases/latest) and move the binary to your `~/.local/bin` directory.


---

## 🚀 Usage

### View Available Options

To view all available commands and options, run:

```bash
wallhaven
```

### Edit Configuration

You can customize the settings, such as specifying the folder where downloaded wallpapers will be saved. To edit the configuration, run:

```bash
wallhaven edit
```

This will open the configuration file, where you can set your preferred download directory.

### Search Wallpapers

To search for wallpapers on Wallhaven, use the search command followed by your query:

```bash
wallhaven search [query]
```

For example:

```bash
wallhaven search mashiro shiina
```

### Download Wallpapers

To download a specific wallpaper by its Wallhaven ID, use the following command:

```bash
wallhaven download [id]
```

For example, to download a wallpaper with the ID `l35l5l`:

```bash
wallhaven download l35l5l
```

> [!Tip]
> You can download multiple images at once:
> 
> ```bash
> wallhaven download l35l5l zy759g
> ```
### Download Collections

To download wallpapers from a user’s collection, use the following command:

```bash
wallhaven collection [username]
```

> [!Tip]
> To download all wallpapers from the collection, add the `-a` flag:
> 
> ```bash
> wallhaven collection [username] -a
> ```

---

### Notes

- Currently, `kitty` is mandatory as the terminal emulator, but future versions may support other terminals.
