#!/bin/bash
set -e

# Variables
#!/bin/bash
set -e

if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <username> <server_ip> <directory1> <directory2> ..."
    exit 1
fi

USERNAME="$1"
SERVER_IP="$2"
shift 2
DIRECTORIES=("$@")

# Function to update the repositories and run Docker Compose
update_directory() {
    local dir=$1
    echo "Updating $dir ..."
    cd $dir
    git fetch
    git pull
    docker-compose -f docker-compose.prod.yml up -d --build
    cd ..
}

# SSH command to connect to the server and update the repositories
ssh -o "StrictHostKeyChecking=no" $USERNAME@$SERVER_IP << 'EOF'
    set -e

    DIRECTORIES=("${@:3}")

    # Iterate through the directories and call the update_directory function
    for dir in "${DIRECTORIES[@]}"; do
        update_directory $dir
    done
EOF

