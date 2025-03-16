#!/bin/bash
# Script to fix remaining test criteria issues

echo "Fixing remaining test criteria issues..."

# Fix 1: Correct the pattern detection tests in legacy pattern_test.go
echo "Fixing pattern criteria in legacy tests (comprehensive fix)"
sed -i.bak '
# Fix float values without explicit type conversion
s/"maxDeviation": \([0-9]\+\.[0-9]\+\)/"maxDeviation": float64(\1)/g
s/"minIncreasePct": \([0-9]\+\.[0-9]\+\)/"minIncreasePct": float64(\1)/g
s/"maxDecreasePct": \([0-9]\+\.[0-9]\+\)/"maxDecreasePct": float64(\1)/g
s/"deviation": \([0-9]\+\.[0-9]\+\)/"deviation": float64(\1)/g
s/"increase": \([0-9]\+\.[0-9]\+\)/"increase": float64(\1)/g
s/"decrease": \([0-9]\+\.[0-9]\+\)/"decrease": float64(\1)/g

# Fix the criteria structure (add $pattern wrapper)
s/criteria := map\[string\]interface{}{\s*\n\s*"pattern":/criteria := map[string]interface{}{\n\t\t\t"$pattern":/g
s/\(maxDeviation.*\),\s*$/\1\n\t\t\t},/g
s/\(minIncreasePct.*\),\s*$/\1\n\t\t\t},/g
s/\(maxDecreasePct.*\),\s*$/\1\n\t\t\t},/g
' tests/legacy/engine/pattern_criteria/pattern_test.go

# Fix 2: Double-check the time_evaluators_test.go file criteria
echo "Double-checking time evaluators test file"
sed -i.bak '
# Fix any remaining float values
s/"maxDeviation": \([0-9]\+\.[0-9]\+\)/"maxDeviation": float64(\1)/g
s/"minIncreasePct": \([0-9]\+\.[0-9]\+\)/"minIncreasePct": float64(\1)/g
s/"maxDecreasePct": \([0-9]\+\.[0-9]\+\)/"maxDecreasePct": float64(\1)/g

# Fix the criteria structure where needed
/consistent/, /}, *$/ {
  s/criteria := map\[string\]interface{}{\s*$/criteria := map[string]interface{}{/
  s/"pattern":\s*"consistent"/"$pattern": map[string]interface{}{\n\t\t"pattern":\t"consistent"/
  s/"maxDeviation": float64\([^)]*\),\s*$/"maxDeviation": float64\1\n\t\t},/
}
' internal/engine/time_evaluators_test.go

# Fix 3: Update the rule_engine_test.go file to follow the pattern
# This is just to check if the criteria nesting is intentional
echo "Updating rule engine test file (criteria nesting)"
sed -i.bak '
# Just make sure all numbers have explicit type conversion
s/"\$[a-z]\+": \([0-9]\+\),/"\$[a-z]\+": float64(\1),/g
' internal/engine/rule_engine_test.go

echo "Script completed. The following files have been potentially modified:"
echo "1. tests/legacy/engine/pattern_criteria/pattern_test.go"
echo "2. internal/engine/time_evaluators_test.go"
echo "3. internal/engine/rule_engine_test.go"

echo "Backup files with .bak extension have been created for modified files."
echo "To verify the fixes, please run:"
echo "  ./tools/validate_test_criteria.sh" 
