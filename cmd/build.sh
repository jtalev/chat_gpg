#!/bin/bash

# Generate secure keys
HMAC_KEY=$(openssl rand -base64 32)
ENC_KEY=$(openssl rand -base64 32)

# File to update
ENV_FILE=".env"

# Ensure the .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "$ENV_FILE not found. Creating a new one."
    touch "$ENV_FILE"
fi

# Update or insert SESSION_HMAC_KEY
if grep -q "^SESSION_HMAC_KEY=" "$ENV_FILE"; then
    # Update the existing SESSION_HMAC_KEY on its own line
    sed -i'' -E "s|^SESSION_HMAC_KEY=.*|SESSION_HMAC_KEY=$HMAC_KEY|" "$ENV_FILE"
else
    # Add SESSION_HMAC_KEY on a new line after the ENV line
    awk -v key="SESSION_HMAC_KEY=$HMAC_KEY" '
        BEGIN {found=0}
        {
            print
            if (!found && /^ENV=/) { print key; found=1 }
        }
        END {if (!found) print key}
    ' "$ENV_FILE" > temp.env && mv temp.env "$ENV_FILE"
fi

# Update or insert SESSION_ENC_KEY
if grep -q "^SESSION_ENC_KEY=" "$ENV_FILE"; then
    # Update the existing SESSION_ENC_KEY on its own line
    sed -i'' -E "s|^SESSION_ENC_KEY=.*|SESSION_ENC_KEY=$ENC_KEY|" "$ENV_FILE"
else
    # Add SESSION_ENC_KEY on a new line after the ENV line
    awk -v key="SESSION_ENC_KEY=$ENC_KEY" '
        BEGIN {found=0}
        {
            print
            if (!found && /^ENV=/) { print key; found=1 }
        }
        END {if (!found) print key}
    ' "$ENV_FILE" > temp.env && mv temp.env "$ENV_FILE"
fi

echo "Keys updated successfully!"

air
