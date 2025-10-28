#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    echo "📄 Loading environment variables from .env..."
    export $(grep -v '^#' .env | xargs)
    echo "✅ Environment variables loaded!"
else
    echo "⚠️  No .env file found. Copy env.example to .env first!"
    exit 1
fi

# Check if OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ] || [ "$OPENAI_API_KEY" = "sk-your-api-key-here" ]; then
    echo "⚠️  OPENAI_API_KEY not configured in .env file"
    echo "   Backend will use mock responses instead of real AI"
else
    echo "✅ OpenAI API Key configured (${OPENAI_API_KEY:0:7}...)"
fi

echo ""
echo "🚀 Starting backend server..."
echo ""

# Run the backend
make run

