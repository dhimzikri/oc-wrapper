@echo off
ren C:\Users\08023.dimas\.bun\bin\opencode.exe opencode.exe.real
copy C:\Users\08023.dimas\.config\opencode\opencode-wrapper\opencode.exe C:\Users\08023.dimas\.bun\bin\opencode.exe
echo Done. Open opencode in new pwsh window to test.
