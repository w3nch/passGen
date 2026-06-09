# Go Password Generator

A simple password generator written in Go.

Features:

* Weak, Medium, and Strong password profiles
* Cryptographically secure password generation using `crypto/rand`
* Mouse movement and system time entropy collection
* SHA-256 entropy hashing
* Terminal menu using `go-pretty`

## Installation

```bash
go mod tidy
go run .
```

## Usage

Select a password strength:

* Weak
* Medium
* Strong

The program will collect entropy, generate a password, and display it.

## Dependencies

* github.com/go-vgo/robotgo
* github.com/jedib0t/go-pretty/v6/table
