# SRR - Surgical RAR Recovery

## Overview

**Surgical RAR Recovery (SRR)** is a high-performance, command-line utility built in Go, designed specifically for recovering passwords of RAR5 archives. Instead of relying on brute force blindly, SRR is tailored for scenarios where a portion of the password is known (masking) and a specific character set is provided.

By leveraging Go's robust concurrency model and avoiding external dependencies, SRR performs deeply efficient, CPU-bound operations strictly in-memory. It surgically parses the RAR5 binary headers to extract the exact cryptographic metadata required for password validation without scanning the entire file.

## Key Features

* **Surgical Extraction:** Directly parses RAR5 binary block headers to retrieve PBKDF2 parameters (salt, iteration count, and password check block) instantly, avoiding full archive loading.
* **Targeted Masking:** Combine known password segments with dynamic character sets to exponentially reduce the search space.
* **High Concurrency:** Utilizes Go's goroutines to distribute cryptographic workload across all available CPU cores, maximizing hash rates.
* **Resilience & State Persistence:** (Planned) Automatically saves progress to a JSON state file. Executions can be paused and resumed exactly where they left off.
* **Zero Dependencies:** Built entirely using the Go standard library (`crypto/sha256`, `crypto/pbkdf2`, etc.), ensuring extreme portability and security.
* **Clean Architecture:** Structured into `domain`, `services`, and `infrastructure` layers, strictly separating business logic from I/O and external adaptations.

## Installation

Ensure you have [Go 1.26+](https://golang.org/dl/) installed.

```bash
git clone https://github.com/UrielJaloto/surgical-rar-recovery.git
cd surgical-rar-recovery
go build -o srr internal/cmd/srr/main.go

```

## Usage

SRR requires the path to the target `.rar` file and a text file containing the allowed character set.

```bash
./srr --file target.rar --charset chars.txt [OPTIONS]

```

### Command-Line Arguments

| Flag | Description | Default | Required |
| --- | --- | --- | --- |
| `--file` | Path to the encrypted RAR5 file | `""` | **Yes** |
| `--charset` | Path to the text file containing the character set | `""` | **Yes** |
| `--known-part` | The known portion of the password | `""` | No |
| `--max-length` | Maximum total password length to attempt | `13` | No |
| `--workers` | Number of concurrent workers (threads) | `1` | No |
| `--state-file` | Path to save/load the execution state | `./state-file.json` | No |

### Example

```bash
./srr --file archive.rar --charset lowercase_alphanumeric.txt --known-part "admin_" --max-length 10 --workers 8

```

## Architecture

The project follows Clean Architecture principles:

* **`internal/domain`**: Contains the core entities (`Config`, `EncryptionMetadata`, `ValidationReport`) and fundamental business rules (e.g., combination math).
* **`internal/services`**: Defines the interfaces and use cases (`RecoveryEngine`, `Application`, `SettingsBuilder`). Connects the domain logic with the infrastructure.
* **`internal/infrastructure`**: Implements the technical details. Contains the RAR5 binary parser, the PBKDF2 cryptographic worker, CLI loaders, and configuration validators.
* **`internal/cmd`**: The entry point that injects the dependencies and orchestrates the application lifecycle.

## Roadmap and Current Status

* [x] **Phase 1: Configuration & Validation:** Implementation of CLI flag parsing, file validation, and complex mathematical checks to warn users of computationally unfeasible configurations.
* [x] **Phase 2: RAR5 Binary Parsing:** Development of a custom reader that navigates the variable-integer structures of RAR5 headers to extract the AES-256 / PBKDF2 encryption metadata (IV, Salt, Check Value).
* [x] **Phase 3: Cryptographic Worker:** Implementation of the core hashing engine utilizing `crypto/pbkdf2` and `crypto/sha256` to validate generated candidates against the extracted metadata.
* [] **Phase 4: Candidate Generator Engine:** Implementation of the password generator that efficiently produces combinations based on the charset and known parts, feeding them into worker channels.
* [ ] **Phase 5: State Management:** Implementation of the pause/resume functionality via the `state-file.json`.
* [ ] **Phase 6: Terminal User Interface (TUI):** Real-time, non-blocking terminal updates displaying current hash rate, progress percentage, and Estimated Time of Arrival (ETA).

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](https://github.com/UrielJaloto/surgical-rar-recovery/blob/develop/LICENSE) file for more details.