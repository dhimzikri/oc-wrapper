# OpenCode Wrapper

Cross-platform wrapper for [OpenCode](https://github.com/anomalyco/opencode) that fixes terminal issues.

## What it fixes

1. **Ctrl+C leak** - Prevents accidental TUI exit when pressing Ctrl+C
2. **Mouse tracking leak** - Cleans up terminal escape sequences on exit (no more weird characters when moving mouse after exit)

## Requirements

- Windows 10/11 or Linux (Ubuntu/Debian/etc.)
- Go 1.21+ (for building, or use pre-built binary)
- OpenCode installed

---

## Windows Installation

### Option 1: Quick Install (PowerShell 7)

```powershell
irm https://raw.githubusercontent.com/dhimzikri/oc-wrapper/master/install.ps1 | iex
```

Then reload your profile:
```powershell
. $PROFILE
```

### Option 2: Manual Install

1. Clone and build:
```powershell
git clone https://github.com/dhimzikri/oc-wrapper.git "$env:USERPROFILE\.config\opencode\wrapper"
cd "$env:USERPROFILE\.config\opencode\wrapper"
go build -o opencode.exe main.go
```

2. Add to PowerShell profile (`$PROFILE`):
```powershell
function opencode { & "$env:USERPROFILE\.config\opencode\wrapper\opencode.exe" @args }
```

3. Reload profile:
```powershell
. $PROFILE
```

---

## Linux Installation

### Option 1: Quick Install

```bash
curl -sSL https://raw.githubusercontent.com/dhimzikri/oc-wrapper/master/install.sh | bash
```

Then reload your shell:
```bash
source ~/.bashrc   # for bash
source ~/.zshrc    # for zsh
```

### Option 2: Manual Install

1. Clone and build:
```bash
git clone https://github.com/dhimzikri/oc-wrapper.git ~/.config/opencode/wrapper
cd ~/.config/opencode/wrapper
GOOS=linux GOARCH=amd64 go build -o opencode main_linux.go
```

2. Add to shell config (`~/.bashrc` or `~/.zshrc`):
```bash
# OpenCode Wrapper
alias opencode='~/.config/opencode/wrapper/opencode'
```

3. Reload config:
```bash
source ~/.bashrc
```

---

## What the installer does

### Windows (install.ps1)

1. Creates `%USERPROFILE%\.config\opencode\wrapper\`
2. Builds or downloads `opencode.exe`
3. Appends function to `$PROFILE`

### Linux (install.sh)

1. Creates `~/.config/opencode/wrapper/`
2. Builds or downloads `opencode`
3. Appends alias to `~/.bashrc` and/or `~/.zshrc`

### What is NOT modified:
- Existing shell config content (only appends, never replaces)
- Your OpenCode installation (`~/.bun/bin/opencode`)
- System PATH
- Other environment variables

---

## Uninstall

### Windows

```powershell
# Remove wrapper files
Remove-Item -Recurse -Force "$env:USERPROFILE\.config\opencode\wrapper"

# Edit profile to remove the function
notepad $PROFILE
# Delete these lines:
#   # OpenCode Wrapper
#   function opencode { & "$env:USERPROFILE\.config\opencode\wrapper\opencode.exe" @args }
```

### Linux

```bash
# Remove wrapper files
rm -rf ~/.config/opencode/wrapper

# Edit shell config to remove the alias
nano ~/.bashrc  # or ~/.zshrc
# Delete these lines:
#   # OpenCode Wrapper
#   alias opencode='~/.config/opencode/wrapper/opencode'
```

---

## Configuration

The wrapper auto-detects opencode in this order:

### Windows
1. `OPENCODE_PATH` environment variable
2. `opencode-core.exe` next to wrapper
3. `~/.bun/bin/opencode.exe`
4. `~/.local/bin/opencode.exe`
5. `~/AppData/Local/Programs/opencode/opencode.exe`
6. System PATH

### Linux
1. `OPENCODE_PATH` environment variable
2. `opencode-core` next to wrapper
3. `~/.bun/bin/opencode`
4. `~/.local/bin/opencode`
5. `/usr/local/bin/opencode`
6. `/usr/bin/opencode`
7. System PATH

To override, set `OPENCODE_PATH`:
```bash
# Linux
export OPENCODE_PATH="/custom/path/opencode"

# Windows PowerShell
$env:OPENCODE_PATH = "C:\custom\path\opencode.exe"
```

---

## Troubleshooting

### "opencode not found" error
The wrapper can't find OpenCode. Either:
1. Install OpenCode first: `bun install -g opencode`
2. Or set the path manually with `OPENCODE_PATH`

### Wrapper installed but `opencode` command not found

**Windows:**
```powershell
. $PROFILE
```

**Linux:**
```bash
source ~/.bashrc
```

Or restart your terminal.

---

## How it works

### Windows
- Uses `SetConsoleCtrlHandler` to block Ctrl+C from reaching the TUI
- Creates child process in separate process group (`CREATE_NEW_PROCESS_GROUP`)
- Resets terminal mouse tracking modes on exit

### Linux
- Uses `signal.Ignore(syscall.SIGINT)` to block Ctrl+C
- Creates child process in separate process group (`Setpgid: true`)
- Resets terminal mouse tracking modes on exit

---

## Building

### Windows (from Windows)
```powershell
go build -o opencode.exe main.go
```

### Linux (from Linux)
```bash
go build -o opencode main_linux.go
```

### Cross-compile Linux from Windows
```powershell
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o opencode-wrapper-linux-amd64 main_linux.go
```

### Cross-compile Windows from Linux
```bash
GOOS=windows GOARCH=amd64 go build -o opencode-wrapper-windows-amd64.exe main.go
```

---

## License

MIT
