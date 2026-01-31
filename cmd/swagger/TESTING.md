# Testing Guide

This README explains how the launch configuration simulates real-world usage by changing the working directory, which is crucial for testing path resolution and file operations in the CLI tool.This README explains how the launch configuration simulates real-world usage by changing the working directory, which is crucial for testing path resolution and file operations in the CLI tool.

## Debugging Configuration

### Launch.json Setup

The project includes a VS Code launch configuration in [.vscode/launch.json](.vscode/launch.json) that allows you to debug the swagger CLI tool while simulating execution from a different directory.

#### Key Configuration Details:

- **Program**: Points to [`cmd/swagger/swagger.go`](cmd/swagger/swagger.go) - the main entry point
- **Working Directory (`cwd`)**: Set to `<YOUR TESTING DIRECTORY>`
- **Arguments**: Configured to run the config command with API validation

#### Why Change Working Directory?

The `cwd` setting simulates running the swagger CLI from a different directory than the source code location. This is important because:

1. **Real-world Usage**: Users will run the CLI from their project directories, not from the swagger tool's source directory
2. **Path Resolution**: Tests relative path handling and file resolution from different working directories  
3. **File Operations**: Validates that file operations work correctly regardless of where the command is executed from
4. **GO MOD INIT** Is a requirement for the project to generate the imports correctly for the .go files
#### Current Test Configuration:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Run Swagger Command",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/swagger/swagger.go",
            "cwd": "absolutePathToTheTestingFolder",
            "args": [
                "config",
                "--api", "{workspaceFolder}/cmd/swagger/commands/testdata", "test_config_api.yaml"
            ]
        }
    ]
}
```

# From the project root
go run cmd/swagger/swagger.go config --api /path/to/testdata test_config_api.yaml

