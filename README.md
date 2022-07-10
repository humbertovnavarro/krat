# Tor Reverse Shell
A reverse shell written in go. Uses the tor network to anonymize remote access, and v3 onion addresses to secure access to remote clients.
Part of a series of tools designed for covert penetration testing.

# Why http?
It's cross platform because go is cross platform. HTTP traffic won't arouse suspicion on the part of the host OS, The complexity of invoking commands is moved to the master node, and it's easier to consume os statistics through a REST api as opposed to ssh.

## Features
- Execute executables on a remote system with args via the tor network securely using a REST api
- Easily command the remote host to download a file using a REST endpoint
- Easily command the remote host to upload a file using a REST endpoint
