#!/bin/bash
echo "Compiling Linux version..."
env GOOS=linux go build

echo "Creating tarball..."
tar -czvf hectorcorrea.tar.gz \
  hectorcorrea.com \
  views/* \
  public/*

echo "Copying to the remote server"
scp hectorcorrea.tar.gz remote_user@remote_server:.

echo "Done."

# From here, login to the remote server, stop the running service,
# untar the file (tar -xzvf hectorcorrea.tar.gz -C target_folder/)
# and restart the service.
