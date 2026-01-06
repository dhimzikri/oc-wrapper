# OpenCode Wrapper

Windows wrapper for [OpenCode](https://github.com/anomalyco/opencode) that fixes terminal issues on Windows/PowerShell.

## What it fixes

1. **Ctrl+C leak** - Prevents accidental TUI exit when pressing Ctrl+C
2. **Mouse tracking leak** - Cleans up terminal escape sequences on exit (no more weird characters when moving mouse after exit)

## Requirements

- Windows 10/11
- Go 1.21+ (for building, or use pre-built binary)
- OpenCode installed

## Installation

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

## What the installer does

The quick install script performs these actions:

### 1. Creates installation directory
```
%USERPROFILE%\.config\opencode\wrapper\
```

### 2. Builds or downloads the wrapper
- **If Go is installed** → Clones repo to temp, builds `opencode.exe`, cleans up
- **If Go is NOT installed** → Downloads pre-built binary from GitHub releases

### 3. Modifies PowerShell profile
**APPENDS** these 2 lines to the END of your `$PROFILE`:
```powershell
# OpenCode Wrapper
function opencode { & "$env:USERPROFILE\.config\opencode\wrapper\opencode.exe" @args }
```

### What is NOT modified:
- Existing PowerShell profile content (only appends, never replaces)
- Your OpenCode installation (`~/.bun/bin/opencode.exe`)
- System PATH
- Other environment variables
- Windows PowerShell 5.x profile (only affects PowerShell 7 if you run from pwsh)

---

## Security & Safety

### Before installing
The installer only modifies your **user-level** PowerShell profile. It does NOT:
- Require administrator privileges
- Modify system files
- Change global PATH
- Affect other users on the machine

### Backup your profile first (recommended)
```powershell
Copy-Item $PROFILE "$PROFILE.backup"
```

---

## Uninstall / Rollback

### Option 1: Automatic uninstall
```powershell
# Remove wrapper files
Remove-Item -Recurse -Force "$env:USERPROFILE\.config\opencode\wrapper"

# Edit profile to remove the function
notepad $PROFILE
# Delete these lines:
#   # OpenCode Wrapper
#   function opencode { & "$env:USERPROFILE\.config\opencode\wrapper\opencode.exe" @args }
```

### Option 2: Restore from backup
If you backed up your profile before installing:
```powershell
Copy-Item "$PROFILE.backup" $PROFILE -Force
. $PROFILE
```

### Option 3: Full profile reset (nuclear option)
If something goes wrong and PowerShell won't start:

1. Open `cmd.exe` (not PowerShell)
2. Delete the profile:
```cmd
del "%USERPROFILE%\Documents\PowerShell\Microsoft.PowerShell_profile.ps1"
```
3. Restart PowerShell - it will start fresh without any profile

---

## Troubleshooting

### PowerShell won't start after install
Your profile has a syntax error. Use the "Full profile reset" option above, or edit the profile from cmd.exe:
```cmd
notepad "%USERPROFILE%\Documents\PowerShell\Microsoft.PowerShell_profile.ps1"
```

### "opencode.exe not found" error
The wrapper can't find OpenCode. Either:
1. Install OpenCode first: `bun install -g opencode`
2. Or set the path manually:
```powershell
$env:OPENCODE_PATH = "C:\path\to\opencode.exe"
```

### Wrapper installed but `opencode` command not found
Reload your profile:
```powershell
. $PROFILE
```
Or restart PowerShell.

---

## Configuration

The wrapper auto-detects opencode.exe in this order:

1. `OPENCODE_PATH` environment variable
2. `opencode-core.exe` next to wrapper (bundled install)
3. `~/.bun/bin/opencode.exe` (bun install)
4. `~/.local/bin/opencode.exe`
5. `~/AppData/Local/Programs/opencode/opencode.exe`
6. System PATH

To override, set `OPENCODE_PATH`:
```powershell
$env:OPENCODE_PATH = "C:\custom\path\opencode.exe"
```

---

## How it works

- Uses Windows `SetConsoleCtrlHandler` to block Ctrl+C from reaching the TUI
- Creates child process in separate process group (`CREATE_NEW_PROCESS_GROUP`)
- Resets terminal mouse tracking modes on exit (prevents character leak)

---

## License

MIT
