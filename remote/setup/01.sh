#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# Set timezone for the server.
TIMEZONE=America/Chicago

# Set name of new user to create.
USERNAME=jambuster

# Prompt to enter a password for PostgreSQL jambuster user.
read -p "Enter password for jambuster DB user: " DB_PASSWORD

# Force all output to be presented in en_US for duration of script.
export LC_ALL=en_US.UTF-8

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

# Enable "universe" repository.
add-apt-repository --yes universe

# Update all software packages.
apt update

# Set the system timezone and install all locales.
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

# Add the new user and give them sudo privileges.
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

# Force password set for new user on initial login.
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

# Copy SSH keys from root user to new user.
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

# Configure the firewall to allow SSH, HTTP, and HTTPS traffic.
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# Install fail2ban.
apt --yes install fail2ban

# Install the migrate CLI tool.
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
mv migrate /usr/local/bin/migrate

# Install PostgreSQL.
apt --yes install postgresql

# Set up the jambuster DB and create a user account with the password entered earlier.
sudo -i -u postgres psql -c "CREATE DATABASE jambuster"
sudo -i -u postgres psql -d jambuster -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d jambuster -c "CREATE ROLE jambuster WITH LOGIN PASSWORD '${DB_PASSWORD}'"

# Add a DSN for connecting to the jambuster database.
echo "JAMBUSTER_DB_DSN='postgres://jambuster:${DB_PASSWORD}@localhost/jambuster'" >> /etc/environment

# Install Caddy.
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt --yes install caddy

# Upgrade all packages.
sudo apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo "Script complete! Rebooting..."
reboot
