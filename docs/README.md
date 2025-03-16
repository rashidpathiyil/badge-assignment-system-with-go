# Badge Assignment System Documentation

This directory contains comprehensive documentation for the Badge Assignment System, with a focus on badge criteria formats and requirements.

## Documentation Index

1. **[Badge Criteria Format Specification](BADGE_CRITERIA_FORMAT.md)** - Complete, detailed specification of all badge criteria formats and requirements.

2. **[Quick Reference Guide](QUICK_REFERENCE_CRITERIA.md)** - Concise reference for common badge criteria formats and API requirements.

3. **[Visual Diagrams](BADGE_CRITERIA_DIAGRAM.md)** - Visual representations of badge criteria structures and flows.

4. **[Examples](examples/)** - Code examples showing properly formatted badge criteria and API requests.

5. **[API Documentation](API.md)** - Detailed API endpoints and usage
6. **[Database Schema](DATABASE_SCHEMA.md)** - Database tables and relationships
7. **[Event Schema](EVENT_SCHEMA.md)** - Event structure and validation
8. **[Badge Criteria](BADGE_CRITERIA.md)** - How to define badge award criteria
9. **[Count Operators](COUNT_OPERATORS.md)** - Operators for counting events
10. **[Time-Based Criteria](TIME_BASED_CRITERIA.md)** - Implementing and testing time-based badge criteria
11. **[Deployment Guide](DEPLOYMENT.md)** - How to deploy the system
12. **[Development Guide](DEVELOPMENT.md)** - Setting up a development environment

## Getting Started

If you're new to the badge assignment system, we recommend starting with:

1. First read the **[Quick Reference Guide](QUICK_REFERENCE_CRITERIA.md)** to understand the basics
2. Review the **[Visual Diagrams](BADGE_CRITERIA_DIAGRAM.md)** to visualize the structures
3. Look at the **[Examples](examples/badge_criteria_examples.go)** for practical implementations
4. Consult the **[Full Specification](BADGE_CRITERIA_FORMAT.md)** for in-depth details

## Common Requirements

- Badge creation requires a `flow_definition` field containing the criteria
- Event submission requires both `data` and `payload` fields
- All numeric values in Go code must use explicit `float64()` conversion
- Pattern criteria must use a `$pattern` wrapper
- Event names must match existing event types exactly

## Validation

Run the validation script to check your criteria format:

```bash
./tools/validate_test_criteria.sh
``` 
