# prosqlctl
This is the installer for prosql-agent

# Notes for Windows
on Windows prosqlctl uses nssm.exe for registering prosql-agent as service.
nssm.exe is sourced from https://nssm.cc/download

# Building
Mac: go build -tags mac

Linux: go build -tags linux

Windows: go build -tags windows

For cross compiling use GOOS as follows:
 
GOOS=windows go build -tags windows
