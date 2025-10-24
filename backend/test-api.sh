#!/bin/bash

# Test script for Stori Backend API
# Usage: ./test-api.sh

BASE_URL="http://localhost:8080"

echo "🧪 Testing Stori Backend API..."
echo "================================"
echo ""

# Test 1: Health Check
echo "1️⃣  Testing Health Endpoint..."
curl -s "$BASE_URL/api/health" | jq '.'
echo ""

# Test 2: API Info
echo "2️⃣  Testing Root Endpoint..."
curl -s "$BASE_URL/" | jq '.'
echo ""

# Test 3: Get All Transactions
echo "3️⃣  Testing Transactions Endpoint..."
curl -s "$BASE_URL/api/transactions" | jq '.count, .period'
echo ""

# Test 4: Get Category Summary
echo "4️⃣  Testing Category Summary..."
curl -s "$BASE_URL/api/summary/categories" | jq '.summary'
echo ""

# Test 5: Get Timeline
echo "5️⃣  Testing Timeline..."
curl -s "$BASE_URL/api/summary/timeline" | jq '.aggregation, .timeline | length'
echo ""

# Test 6: Date Range Filter
echo "6️⃣  Testing Date Range Filter..."
curl -s "$BASE_URL/api/transactions?startDate=2024-01-01&endDate=2024-01-31" | jq '.count'
echo ""

echo "✅ All tests complete!"


