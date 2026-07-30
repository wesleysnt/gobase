#!/usr/bin/env bash
set -euo pipefail

echo "GoBase Project Setup"
echo "===================="
echo ""

# Prompt for new module path
read -r -p "Module path (e.g., github.com/you/myproject): " MODULE_PATH
if [ -z "$MODULE_PATH" ]; then
    echo "Error: Module path is required"
    exit 1
fi

OLD_MODULE="github.com/you/gobase"

echo ""
echo "Replacing module references..."
echo "  $OLD_MODULE → $MODULE_PATH"

# Replace in go.mod
if [ -f go.mod ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|${OLD_MODULE}|${MODULE_PATH}|g" go.mod
    else
        sed -i "s|${OLD_MODULE}|${MODULE_PATH}|g" go.mod
    fi
fi

# Replace in all .go files
find . -name "*.go" -type f | while read -r file; do
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|${OLD_MODULE}|${MODULE_PATH}|g" "$file"
    else
        sed -i "s|${OLD_MODULE}|${MODULE_PATH}|g" "$file"
    fi
done

echo "Module path set to: $MODULE_PATH"

# Copy .env.example to .env
if [ -f .env.example ] && [ ! -f .env ]; then
    cp .env.example .env
    echo ".env created from .env.example — edit with your settings"
elif [ -f .env ]; then
    echo ".env already exists, skipping"
fi

# Run go mod tidy
echo ""
echo "Running go mod tidy..."
go mod tidy

echo ""
echo "===================="
echo "Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Edit .env with your database URL and JWT secret"
echo "  2. Run: make dev"
