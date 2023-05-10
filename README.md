# Script-utils
This repository contain usual script needed to install package or configure server

## The `setup_ssh_key_auth.sh` script
The script takes these values as command-line arguments in the following format: 
```bash
./setup_ssh_key_auth.sh <username> <server_ip> <email>.
```

Replace <username> with your username, <server_ip> with your server's IP address or hostname, and <email> with your email address. The script will generate an SSH key pair if it doesn't exist and copy the public key to the remote server. Finally, it will test the SSH key-based authentication and display a success message if everything is set up correctly.

## The `update_directory.sh` script
The script now takes the following arguments:
```bash
./update_directory.sh <username> <server_ip> <directory1> <directory2> ...
```
This script connects to the server using SSH key-based authentication, navigates to the directories provided as arguments , performs a git fetch, git pull, and runs docker-compose with the specified production file.

## The `install_docker.sh` script
This script installs Docker and Docker Compose on your Ubuntu server, enables and starts the Docker service, and prints the installed versions of Docker and Docker Compose. If you want to install a different version of Docker Compose, replace 1.29.2 in the script with the desired version.