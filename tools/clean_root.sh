#!/bin/bash

# Script to clean up the root directory
echo "Cleaning up root directory..."

# Create necessary directories
mkdir -p scripts
mkdir -p bin
mkdir -p reports

# Move shell scripts to scripts directory
echo "Moving shell scripts to scripts directory..."
for script in *.sh; do
    if [ -f "$script" ]; then
        echo "Moving $script to scripts/"
        mv "$script" "scripts/"
    fi
done

# Move binary files to bin directory
echo "Moving binary files to bin directory..."
if [ -f "server" ]; then
    echo "Moving server to bin/"
    mv "server" "bin/"
fi

if [ -f "badgecli" ]; then
    echo "Moving badgecli to bin/"
    mv "badgecli" "bin/"
fi

# Move log files to logs directory
echo "Moving log files to logs directory..."
if [ -f "server_logs.txt" ]; then
    echo "Moving server_logs.txt to tests/logs/"
    mv "server_logs.txt" "tests/logs/"
fi

# Move coverage report to reports directory
echo "Moving coverage report to reports directory..."
if [ -f "coverage.out" ]; then
    echo "Moving coverage.out to reports/"
    mv "coverage.out" "reports/"
fi

# Create README files for new directories
echo "Creating README files for new directories..."

# Create README for scripts directory
cat > "scripts/README.md" << 'EOF'
# Scripts Directory

This directory contains utility scripts for the badge assignment system.

## Available Scripts

These scripts provide shortcuts for common development and testing tasks.
EOF

# Create README for bin directory
cat > "bin/README.md" << 'EOF'
# Binary Directory

This directory contains compiled binaries for the badge assignment system.

**Note:** These files should not be committed to version control.
EOF

# Create README for reports directory
cat > "reports/README.md" << 'EOF'
# Reports Directory

This directory contains generated reports from the badge assignment system.

## Available Reports

- `coverage.out`: Test coverage report generated with `go test -coverprofile=coverage.out`

To view coverage as HTML:

```bash
go tool cover -html=reports/coverage.out
```

**Note:** These files should not be committed to version control.
EOF

# Update .gitignore to include new directories
echo "Updating .gitignore..."
cat >> ".gitignore" << 'EOF'

# Binaries
bin/*
!bin/README.md

# Reports
reports/*
!reports/README.md

# Scripts output
scripts/*.log
EOF

echo "Root directory cleanup complete!"
echo "Files are now organized into appropriate directories." 
