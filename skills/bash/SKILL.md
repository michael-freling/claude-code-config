---
name: bash
description: Bash scripting following best practices for error handling, ShellCheck compliance, testing, code organization, and reliability. Use when adding or updating Bash scripts in projects.
allowed-tools: [Read, Write, Edit, Glob, Grep, Bash]
---

# Bash Scripting Skill

A comprehensive skill for creating and updating Bash scripts following industry best practices from ShellCheck, Google's Shell Style Guide, and 2025 modern shell scripting standards.

## When to Use

Use this skill when:
- Creating new Bash scripts
- Updating or enhancing existing shell scripts
- Need to follow shell scripting best practices
- Working with automation, CI/CD, or build scripts

## Process

### 0. Read Project Design Documentation

**CRITICAL FIRST STEP: Always check for and read `.claude/design.md`**

Before starting any implementation:

1. **Look for `.claude/design.md` in the current directory**
   - If found, read it thoroughly
   - This contains project-specific coding standards, conventions, and architecture
   - Follow these guidelines strictly as they override general best practices

2. **For monorepos or subprojects:**
   - Check for `.claude/design.md` in the subproject root
   - Also check the repository root for overall standards
   - Subproject-specific rules take precedence over repository-level rules

3. **If no design.md exists:**
   - Consider running `/document-design` to create one
   - Or proceed with analyzing the codebase manually

**What to extract from design.md:**
- Project-specific naming conventions
- Script organization patterns
- Error handling conventions
- Testing requirements
- Logging patterns
- Code examples showing preferred style

### 1. Analyze Project Structure

- Check for existing scripts and their patterns
- Identify script organization (scripts/, bin/, tools/)
- Review existing error handling approaches
- Check for CI/CD integration
- Identify testing patterns (bats, shunit2, etc.)

### 2. Search for Relevant Code

- Search for similar script implementations
- Find common utility functions
- Look for error handling patterns
- Identify logging conventions

## Core Principles

### Simplicity First
- **DRY (Don't Repeat Yourself)**: Extract common patterns into functions or shared scripts
- **Prefer Breaking Changes**: Unless explicitly mentioned, prefer simplicity over backward compatibility
- **Early Returns**: Use early returns instead of nested conditionals
- **Minimal Scope**: Keep variables local to functions when possible

### Consistency
- Follow project-specific patterns from `.claude/design.md`
- Match existing script style and naming conventions
- Use consistent error handling patterns throughout

### Fail-Fast
- **CRITICAL: Always handle errors explicitly**
- Use `set -euo pipefail` at the beginning of every script
- Validate inputs early and fail with meaningful error messages
- Never silently swallow errors

### Latest Versions
- Use modern Bash features (Bash 4.0+ when available)
- Prefer built-in features over external commands when possible
- Document minimum Bash version requirements if using newer features

### Comments
- **Comments MUST BE about WHY not WHAT** - Explain the reasoning behind the code, not what the code does
- The code itself should be self-explanatory for what it does
- Use comments to explain complex logic, gotchas, or non-obvious decisions

## Bash Best Practices

### Strict Mode (Error Handling)

**CRITICAL: Every script MUST start with strict error handling**

```bash
#!/usr/bin/env bash

# Strict mode - REQUIRED for all scripts
set -euo pipefail

# Optional: Enable debug mode for development
# set -x

# Example script with proper error handling
```

**What each option does:**
- `set -e` (errexit): Exit immediately if any command fails
- `set -u` (nounset): Treat unset variables as errors
- `set -o pipefail`: Return exit code of first failed command in a pipeline

Example:
```bash
#!/usr/bin/env bash
set -euo pipefail

# Good: Script exits if directory creation fails
mkdir /tmp/myapp
cd /tmp/myapp
echo "Setup complete"

# Good: Explicit error handling with trap
cleanup() {
    local exit_code=$?
    echo "Cleaning up..."
    rm -rf /tmp/myapp
    exit "${exit_code}"
}
trap cleanup EXIT ERR

# Good: Manual error checking when needed
if ! command -v docker &> /dev/null; then
    echo "Error: docker is not installed" >&2
    exit 1
fi
```

### ShellCheck Integration

**CRITICAL: All scripts must pass ShellCheck validation**

ShellCheck is a static analysis tool that catches common bugs and anti-patterns.

```bash
# Install ShellCheck
# Ubuntu/Debian:
sudo apt-get install shellcheck

# macOS:
brew install shellcheck

# Check a script
shellcheck script.sh

# Check all scripts in a directory
shellcheck scripts/*.sh

# Integrate with CI/CD
shellcheck --format=gcc scripts/*.sh
```

**Address all ShellCheck warnings before committing code.**

Example of common issues ShellCheck catches:
```bash
# Bad: Unquoted variable expansion
rm -rf $directory/*

# Good: Quoted to prevent word splitting
rm -rf "${directory:?}"/*

# Bad: Using backticks for command substitution
result=`date +%s`

# Good: Use $() syntax
result=$(date +%s)

# Bad: Ignoring command exit status
grep "pattern" file.txt

# Good: Explicit handling
if grep -q "pattern" file.txt; then
    echo "Pattern found"
fi
```

### Variable Naming and Usage

```bash
# Good: Use descriptive names
readonly CONFIG_FILE="/etc/myapp/config.conf"
readonly MAX_RETRIES=3
user_input=""

# Good: Use UPPER_CASE for environment variables and constants
export DATABASE_URL="postgresql://localhost/mydb"

# Good: Use lower_case for local variables
local temp_file="/tmp/data.txt"

# Good: Always quote variables
echo "User: ${user_name}"
rm -f "${temp_file}"

# Good: Use braces for clarity
echo "File: ${file_name}.txt"

# Good: Use readonly for constants
readonly VERSION="1.0.0"

# Good: Check if variable is set
if [[ -z "${var:-}" ]]; then
    echo "Variable is not set"
fi

# Good: Provide default values
database_host="${DATABASE_HOST:-localhost}"
port="${PORT:-8080}"
```

### Function Design

```bash
# Good: Function documentation and structure
##
# Processes user data and writes to output file
# Arguments:
#   $1 - Input file path
#   $2 - Output file path
# Returns:
#   0 on success, 1 on error
##
process_data() {
    local input_file="${1:?Input file required}"
    local output_file="${2:?Output file required}"

    # Validate inputs early
    if [[ ! -f "${input_file}" ]]; then
        echo "Error: Input file does not exist: ${input_file}" >&2
        return 1
    fi

    # Process data
    while IFS= read -r line; do
        echo "Processed: ${line}" >> "${output_file}"
    done < "${input_file}"

    return 0
}

# Good: Use local variables in functions
calculate_total() {
    local -i sum=0
    local item

    for item in "$@"; do
        ((sum += item))
    done

    echo "${sum}"
}

# Good: Early return pattern
validate_config() {
    local config_file="${1:?Config file required}"

    # Early return for errors
    [[ -f "${config_file}" ]] || {
        echo "Error: Config file not found" >&2
        return 1
    }

    [[ -r "${config_file}" ]] || {
        echo "Error: Config file not readable" >&2
        return 1
    }

    # Main logic only if validation passes
    echo "Config is valid"
    return 0
}
```

### Command Substitution and Pipelines

```bash
# Good: Use $() instead of backticks
current_date=$(date +%Y-%m-%d)
file_count=$(find . -type f | wc -l)

# Good: Proper pipeline error handling with pipefail
set -o pipefail
cat large_file.txt | grep "pattern" | sort | uniq > results.txt

# Good: Process command output line by line
while IFS= read -r line; do
    echo "Processing: ${line}"
done < <(find . -name "*.txt")

# Good: Capture both stdout and stderr
output=$(command 2>&1)
exit_code=$?
```

### Conditionals and Tests

```bash
# Good: Use [[ ]] for tests (Bash built-in, more features)
if [[ -f "${file}" ]]; then
    echo "File exists"
fi

# Good: String comparisons
if [[ "${status}" == "success" ]]; then
    echo "Operation completed successfully"
fi

# Good: Numeric comparisons
if (( count > 10 )); then
    echo "Count exceeds threshold"
fi

# Good: Multiple conditions
if [[ -f "${file}" && -r "${file}" ]]; then
    cat "${file}"
fi

# Good: Pattern matching
if [[ "${filename}" == *.txt ]]; then
    echo "Text file detected"
fi

# Good: Early exit pattern (prefer over nested if)
if [[ ! -f "${config_file}" ]]; then
    echo "Error: Config not found" >&2
    exit 1
fi

if [[ ! -r "${config_file}" ]]; then
    echo "Error: Config not readable" >&2
    exit 1
fi

# Main logic here
```

### Arrays

```bash
# Good: Declare and use arrays
declare -a files=()
files=("file1.txt" "file2.txt" "file3.txt")

# Good: Iterate over array
for file in "${files[@]}"; do
    echo "Processing: ${file}"
done

# Good: Array length
echo "Total files: ${#files[@]}"

# Good: Add to array
files+=("file4.txt")

# Good: Associative arrays (Bash 4.0+)
declare -A config
config[host]="localhost"
config[port]="8080"

for key in "${!config[@]}"; do
    echo "${key}: ${config[${key}]}"
done
```

### Logging and Output

```bash
# Good: Separate stdout and stderr
log_info() {
    echo "[INFO] $*"
}

log_error() {
    echo "[ERROR] $*" >&2
}

log_debug() {
    if [[ "${DEBUG:-0}" == "1" ]]; then
        echo "[DEBUG] $*" >&2
    fi
}

# Usage
log_info "Starting application"
log_error "Failed to connect to database"
DEBUG=1 log_debug "Variable value: ${var}"

# Good: Timestamp logging
log_with_timestamp() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*"
}

# Good: Redirect output to log file
exec 1> >(tee -a "/var/log/myapp.log")
exec 2>&1
```

### Script Organization

```bash
#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# Script: deploy.sh
# Description: Deploys application to production
# Usage: ./deploy.sh [options]
# ============================================================================

# ----------------------------------------------------------------------------
# Constants
# ----------------------------------------------------------------------------
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
readonly VERSION="1.0.0"

# ----------------------------------------------------------------------------
# Global Variables
# ----------------------------------------------------------------------------
DRY_RUN=0
VERBOSE=0

# ----------------------------------------------------------------------------
# Functions
# ----------------------------------------------------------------------------

show_help() {
    cat << EOF
Usage: ${0##*/} [OPTIONS]

Deploy application to production environment.

OPTIONS:
    -h, --help          Show this help message
    -n, --dry-run       Run in dry-run mode
    -v, --verbose       Enable verbose output
    --version           Show version

EXAMPLES:
    ${0##*/}                    # Normal deployment
    ${0##*/} --dry-run          # Test deployment
    ${0##*/} --verbose          # Verbose output
EOF
}

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                show_help
                exit 0
                ;;
            -n|--dry-run)
                DRY_RUN=1
                shift
                ;;
            -v|--verbose)
                VERBOSE=1
                set -x
                shift
                ;;
            --version)
                echo "${VERSION}"
                exit 0
                ;;
            *)
                echo "Error: Unknown option: $1" >&2
                show_help
                exit 1
                ;;
        esac
    done
}

main() {
    parse_arguments "$@"

    log_info "Starting deployment"

    # Main logic here

    log_info "Deployment complete"
}

# ----------------------------------------------------------------------------
# Main
# ----------------------------------------------------------------------------
main "$@"
```

### Testing

**CRITICAL: Use table-driven testing approaches when possible**

Bash testing frameworks provide structure for reliable tests:

1. **bats-core** - Most popular, supports Bash only
2. **shunit2** - Older but stable, supports multiple shells
3. **ShellSpec** - BDD-style, modern features

**Example with bats-core:**

```bash
#!/usr/bin/env bats

# test_functions.bats

setup() {
    # Runs before each test
    export TEST_DIR="/tmp/test_$$"
    mkdir -p "${TEST_DIR}"
}

teardown() {
    # Runs after each test
    rm -rf "${TEST_DIR}"
}

# Good: Test function with multiple cases (table-driven approach)
@test "calculate_sum handles various inputs" {
    # Load functions
    source functions.sh

    # Test cases - happy path
    run calculate_sum 1 2 3
    [ "$status" -eq 0 ]
    [ "$output" -eq 6 ]

    run calculate_sum 0 0 0
    [ "$status" -eq 0 ]
    [ "$output" -eq 0 ]

    run calculate_sum -1 1
    [ "$status" -eq 0 ]
    [ "$output" -eq 0 ]
}

@test "calculate_sum handles error cases" {
    source functions.sh

    # Test error case - no arguments
    run calculate_sum
    [ "$status" -eq 1 ]

    # Test error case - non-numeric input
    run calculate_sum "abc"
    [ "$status" -eq 1 ]
}

@test "validate_config detects missing files" {
    source functions.sh

    run validate_config "/nonexistent/config.conf"
    [ "$status" -eq 1 ]
    [[ "$output" =~ "not found" ]]
}

@test "validate_config accepts valid config" {
    source functions.sh

    local config_file="${TEST_DIR}/config.conf"
    echo "key=value" > "${config_file}"

    run validate_config "${config_file}"
    [ "$status" -eq 0 ]
}
```

**Example with shunit2:**

```bash
#!/bin/bash

# test_functions.sh

# Source the functions to test
. ./functions.sh

# Setup runs before each test
setUp() {
    TEST_DIR="/tmp/test_$$"
    mkdir -p "${TEST_DIR}"
}

# Teardown runs after each test
tearDown() {
    rm -rf "${TEST_DIR}"
}

# Good: Table-driven test approach
test_calculate_sum_happy_cases() {
    # Test case 1
    result=$(calculate_sum 1 2 3)
    assertEquals "Sum of 1 2 3" 6 "${result}"

    # Test case 2
    result=$(calculate_sum 0 0 0)
    assertEquals "Sum of 0 0 0" 0 "${result}"

    # Test case 3
    result=$(calculate_sum -1 1)
    assertEquals "Sum of -1 1" 0 "${result}"
}

test_calculate_sum_error_cases() {
    # Test no arguments
    calculate_sum 2>/dev/null
    assertNotEquals "No arguments should fail" 0 $?

    # Test non-numeric input
    calculate_sum "abc" 2>/dev/null
    assertNotEquals "Non-numeric should fail" 0 $?
}

# Load shunit2
. shunit2
```

**Running tests:**

```bash
# With bats-core
bats test_functions.bats

# With shunit2
./test_functions.sh

# In CI/CD
shellcheck *.sh && bats tests/*.bats
```

### Debugging

```bash
# Good: Enable debug mode
set -x  # Print each command before execution
set +x  # Disable debug mode

# Good: Debug specific section
debug_section() {
    set -x
    # Commands to debug
    complex_operation
    set +x
}

# Good: Conditional debugging
if [[ "${DEBUG:-0}" == "1" ]]; then
    set -x
fi

# Good: Use bash -n for syntax checking
bash -n script.sh

# Good: Use bash -x for trace execution
bash -x script.sh
```

## Implementation Strategy

**Step 1: Review Project Guidelines**
- Read `.claude/design.md` if it exists (MANDATORY)
- Extract project-specific patterns
- Note script organization conventions
- Identify testing requirements

**Step 2: Plan the Changes**
- Create a todo list with specific tasks
- Break down into testable units
- Identify reusable functions

**Step 3: Read Existing Code**
- Review similar script implementations
- Understand existing error handling patterns
- Note utilities and helper functions available

**Step 4: Implement Changes**
- Start with `set -euo pipefail`
- Follow patterns from design.md
- Use naming conventions from project
- Write clear, readable code
- Add appropriate comments for complex logic

**Step 5: Ensure Quality and Fix All Issues**

**CRITICAL: All quality checks must pass before considering the task complete.**

1. **Run syntax check**
   ```bash
   bash -n script.sh
   ```
   - If syntax errors exist, fix them immediately
   - Do not proceed until syntax is valid

2. **Run ShellCheck (MANDATORY)**
   ```bash
   shellcheck script.sh
   ```
   - If ShellCheck reports issues, fix them immediately
   - Address all warnings and errors
   - Do not proceed until ShellCheck passes with no warnings

3. **Test the script**
   ```bash
   # Manual testing
   bash script.sh --help
   bash script.sh --dry-run

   # Automated testing (if available)
   bats tests/test_script.bats
   # or
   ./test_script.sh  # shunit2
   ```
   - If tests fail, fix them immediately
   - Add new tests if implementing new features
   - Update existing tests if modifying behavior
   - Do not proceed until all tests pass

4. **Run in debug mode to verify logic**
   ```bash
   bash -x script.sh
   ```
   - Verify the execution flow
   - Check for unexpected behavior

5. **Check permissions**
   ```bash
   chmod +x script.sh
   ```
   - Ensure script is executable

**Iterative Fix Process:**
- If any step fails, fix the issues and re-run ALL previous steps
- Continue iterating until ALL checks pass:
  - ✅ Syntax valid (bash -n)
  - ✅ No ShellCheck warnings
  - ✅ All tests pass
  - ✅ Script executes correctly
  - ✅ Proper permissions set

## Commands Reference

```bash
# Syntax checking
bash -n script.sh

# ShellCheck
shellcheck script.sh
shellcheck --severity=warning script.sh
shellcheck --format=gcc scripts/*.sh  # CI-friendly

# Testing
bats tests/*.bats                      # bats-core
./test_script.sh                       # shunit2

# Debugging
bash -x script.sh                      # Trace execution
set -x                                 # Enable tracing
set +x                                 # Disable tracing

# Permissions
chmod +x script.sh
chmod 755 script.sh
```

## Common Pitfalls to Avoid

### Unquoted Variables
```bash
# Bad: Word splitting issues
rm -rf $directory

# Good: Always quote
rm -rf "${directory}"

# Bad: Glob expansion issues
for file in $(ls *.txt); do
    echo "${file}"
done

# Good: Use glob directly
for file in *.txt; do
    echo "${file}"
done
```

### Ignoring Exit Codes
```bash
# Bad: Ignoring failures
grep "pattern" file.txt
echo "Done"

# Good: Check exit code
if grep -q "pattern" file.txt; then
    echo "Pattern found"
else
    echo "Pattern not found"
fi

# Good: Use set -e to auto-fail
set -e
critical_command
echo "This won't run if critical_command fails"
```

### Using cd Without Error Handling
```bash
# Bad: cd might fail silently
cd /some/directory
rm -rf *

# Good: Check cd success
cd /some/directory || exit 1
rm -rf *

# Better: Use subshell to auto-return
(
    cd /some/directory || exit 1
    rm -rf *
)
```

### Useless Use of cat
```bash
# Bad: Unnecessary cat
cat file.txt | grep "pattern"

# Good: Direct input redirection
grep "pattern" < file.txt

# Better: Pass filename as argument
grep "pattern" file.txt
```

## Checklist

### Before Starting
- [ ] **Read `.claude/design.md` if it exists** (CRITICAL)
- [ ] Extract project-specific conventions
- [ ] Analyze existing script structure
- [ ] Search for similar implementations
- [ ] Review existing error handling patterns

### During Implementation
- [ ] Add `set -euo pipefail` at script start
- [ ] Use `#!/usr/bin/env bash` shebang
- [ ] Quote all variable expansions
- [ ] Use `[[ ]]` for conditionals
- [ ] Use `$()` instead of backticks
- [ ] Add function documentation
- [ ] Implement early returns for error cases
- [ ] Use local variables in functions
- [ ] Add meaningful error messages to stderr

### After Implementation - MUST ALL PASS
- [ ] Run `bash -n script.sh` - **FIX ALL SYNTAX ERRORS**
- [ ] Run `shellcheck script.sh` - **FIX ALL WARNINGS**
- [ ] Run tests if available - **ALL TESTS MUST PASS**
- [ ] Test script manually with various inputs
- [ ] Run in debug mode (`bash -x`) to verify flow
- [ ] Set executable permissions (`chmod +x`)
- [ ] **Iterate until all checks pass** - do not stop until everything is green
- [ ] Update documentation/comments

## Key Principles

1. **Project Guidelines First**: Always read and follow `.claude/design.md`
2. **Error Handling**: Always use `set -euo pipefail` and handle errors explicitly
3. **ShellCheck**: All scripts must pass ShellCheck validation
4. **Simplicity**: Write clear, maintainable scripts
5. **Quote Everything**: Always quote variable expansions
6. **Testing**: Add automated tests when possible (bats, shunit2)
7. **Fail-Fast**: Validate inputs early and exit on errors

## Version History

- **v1.0.0** (2025-01-16): Initial version with comprehensive Bash best practices
  - Strict mode error handling (set -euo pipefail)
  - ShellCheck integration and compliance
  - Testing frameworks (bats-core, shunit2)
  - Function design and variable naming conventions
  - Debugging and logging patterns
  - Table-driven testing approach
  - Integration with general coding guidelines (DRY, consistency, fail-fast)
