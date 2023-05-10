#!/bin/bash
set -e

# Variables
HOST="root"
SERVER="185.182.187.88"
DIRECTORIES=("mixdata-back" "mixdata-front")

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
ssh -o "StrictHostKeyChecking=no" $HOST@$SERVER << 'EOF'
    set -e

    # Iterate through the directories and call the update_directory function
    for dir in "${DIRECTORIES[@]}"; do
        update_directory $dir
    done
EOF
