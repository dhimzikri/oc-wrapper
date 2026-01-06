# PowerShell function to run opencode through wrapper
# Add this to $PROFILE

function opencode {
    $wrapper = "C:\Users\08023.dimas\.config\opencode\opencode-wrapper\opencode.exe"
    $opencode = "C:\Users\08023.dimas\.bun\bin\opencode.exe"

    if (Test-Path $wrapper) {
        & $wrapper $args
    } else {
        & $opencode $args
    }
}

# To install: Add this line to $PROFILE
# . C:\Users\08023.dimas\.config\opencode\opencode-wrapper\opencode.ps1
