#!/bin/bash

# Final cleanup script to remove all unnecessary files
echo "Performing final cleanup of the codebase..."

# Remove all backup files
echo "Removing backup files..."
find tests -type f -name "*.bak" -delete

# Remove temporary fix scripts that are no longer needed
echo "Removing temporary fix scripts..."
rm -f tools/fix_all_tests.sh
rm -f tools/fix_case.sh
rm -f tools/fix_imports.sh
rm -f tools/fix_imports_final.sh
rm -f tools/fix_imports_properly.sh
rm -f tools/fix_missing_imports.sh
rm -f tools/fix_redeclarations.sh
rm -f tools/fix_unused_imports.sh

# Check if there are any integration_backup directories we can remove
echo "Checking for old integration test backups..."
backup_dirs_count=$(find tests -type d -name "integration_backup_*" | wc -l)
if [ "$backup_dirs_count" -gt 1 ]; then
    echo "Found multiple backup directories. Keeping only the most recent one."
    # Get the most recent backup directory
    most_recent=$(find tests -type d -name "integration_backup_*" | sort | tail -n 1)
    # Remove all other backup directories
    for dir in $(find tests -type d -name "integration_backup_*" | sort | head -n -1); do
        echo "Removing old backup: $dir"
        rm -rf "$dir"
    done
    echo "Kept most recent backup: $most_recent"
else
    echo "Only one backup directory found. Keeping it for reference."
fi

echo "Final cleanup complete! The codebase is now clean and ready for development." 
