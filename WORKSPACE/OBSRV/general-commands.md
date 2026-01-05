#!/bin/bash

if [[ $(whoami) == "sanketika" ]]; then
        echo "script to used by users other than 'sanketika'"
        exit 1
fi

scp /home/sanketika/.kube/config ~/.kube/config

# Function to display usage information
usage() {
    echo "======================"
    echo "Usage: $0 [--kubeconfig <path_to_kubeconfig>] [--proxy <proxy_url>]"
    echo "You can also use KUBECONFIG and PROXY_URL environment variables."
    echo "Command-line arguments take precedence over environment variables."
    echo "Example: $0 --kubeconfig /path/to/kubeconfig --proxy socks5://localhost:1080"
    echo "======================"
    echo ""
}

usage

# Initialize variables with environment variables (if set)
KUBECONFIG=${KUBECONFIG:-""}
PROXY_URL=${PROXY_URL:-""}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --kubeconfig)
            KUBECONFIG="$2"
            shift 2
            ;;
        --proxy)
            PROXY_URL="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown argument: $1"
            usage
            ;;
    esac
done

# Check if kubeconfig file is provided
if [ -z "$KUBECONFIG" ]; then
    echo "Kubeconfig not specified. Using default '~/.kube/config'"
    KUBECONFIG=~/.kube/config
fi


# Check if kubeconfig file exists
if [ ! -f "$KUBECONFIG" ]; then
    echo "Error: Kubeconfig file does not exist: $KUBECONFIG"
    exit 1
fi

echo "Using KUBECONFIG=$KUBECONFIG and PROXY=$PROXY_URL"


# Temporary file to store the modified kubeconfig
TEMP_KUBECONFIG=$(mktemp)

# Function to replace file path with base64 encoded file content
replace_with_base64_content() {
    local path=$1
    if [ -f "$path" ]; then
        content=$(cat "$path" | base64 -w 0)
        echo "    $2: $content"
    else
        echo "    $2: $path  # File not found"
    fi
}

# Process the kubeconfig file
while IFS= read -r line; do
    if [[ $line =~ certificate-authority:\ *(.*) ]]; then
        replace_with_base64_content "${BASH_REMATCH[1]}" "certificate-authority-data"
    elif [[ $line =~ client-certificate:\ *(.*) ]]; then
        replace_with_base64_content "${BASH_REMATCH[1]}" "client-certificate-data"
    elif [[ $line =~ client-key:\ *(.*) ]]; then
        replace_with_base64_content "${BASH_REMATCH[1]}" "client-key-data"
    elif [[ $line =~ ^([[:space:]]*)server:.*$ ]]; then
        echo "$line"
	if [ -n "$PROXY_URL" ]; then
            indentation="${BASH_REMATCH[1]}"
            echo "${indentation}proxy-url: $PROXY_URL"
        fi
    else
        echo "$line"
    fi
done < "$KUBECONFIG" > "$TEMP_KUBECONFIG"

# Replace the original kubeconfig with the modified one
mv "$TEMP_KUBECONFIG" "$KUBECONFIG"

echo "Kubeconfig updated successfully."