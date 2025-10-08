#!/bin/bash

# Perimeter81 Provider Registration Audit Script
# This script helps identify missing resource and data source registrations

echo "=== Perimeter81 Terraform Provider Registration Audit ==="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "ERROR: Must run from repository root (where go.mod is located)"
    exit 1
fi

echo "📁 Repository root: $(pwd)"
echo ""

# Find all resource files
echo "=== Resource Files Found ==="
RESOURCE_FILES=$(ls perimeter81/resource_*.go 2>/dev/null)
if [ -z "$RESOURCE_FILES" ]; then
    echo "❌ No resource files found"
else
    echo "$RESOURCE_FILES" | while read file; do
        # Extract the function name (e.g., resource_network.go -> resourceNetwork)
        basename=$(basename "$file" .go)
        # Convert to camelCase function name
        funcname=$(echo "$basename" | sed 's/_\([a-z]\)/\U\1/g' | sed 's/resource/resource/')
        echo "  ✓ $file"
        echo "    → Function: ${funcname}()"
    done
fi
echo ""

# Find all data source files
echo "=== Data Source Files Found ==="
DATASOURCE_FILES=$(ls perimeter81/data_source_*.go 2>/dev/null)
if [ -z "$DATASOURCE_FILES" ]; then
    echo "❌ No data source files found"
else
    echo "$DATASOURCE_FILES" | while read file; do
        basename=$(basename "$file" .go)
        funcname=$(echo "$basename" | sed 's/_\([a-z]\)/\U\1/g')
        echo "  ✓ $file"
        echo "    → Function: ${funcname}()"
    done
fi
echo ""

# Check provider.go registrations
echo "=== Checking provider.go Registrations ==="
if [ ! -f "perimeter81/provider.go" ]; then
    echo "❌ provider.go not found!"
    exit 1
fi

echo ""
echo "📋 Registered Resources in provider.go:"
grep -A 100 'ResourcesMap.*map\[string\]\*schema.Resource' perimeter81/provider.go | grep 'perimeter81_' | sed 's/^[ \t]*/  /' || echo "  ❌ Could not parse ResourcesMap"

echo ""
echo "📋 Registered Data Sources in provider.go:"
grep -A 100 'DataSourcesMap.*map\[string\]\*schema.Resource' perimeter81/provider.go | grep 'perimeter81_' | sed 's/^[ \t]*/  /' || echo "  ❌ Could not parse DataSourcesMap"

echo ""
echo "=== Recommendations ==="
echo "1. Compare the files found above with the registrations in provider.go"
echo "2. Add any missing resources/data sources to provider.go"
echo "3. Check for copy-paste errors (like resourceObjectServicesRead in resource_object_addresses.go)"
echo ""
