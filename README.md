# todo-cli

A small, dependency-light command-line todo manager written in Go. Tasks are stored in a plain CSV file, so the data stays readable, greppable, and easy to back up — no database, no daemon, no network access.

## Features

- **Add, list, complete, and delete** tasks from the terminal
- **Filtering** — hide completed tasks by default, show everything with `--all`
- **Human-readable timestamps** — creation times are rendered as relative ages ("2 hours ago")
- **Advisory file lock** — `flock` is taken on `todo.csv` before each read-modify-write cycle (see [Known Limitations](#known-limitations) for the gap)
- **Portable data** — a single `todo.csv` file, version-control friendly
- **Unix-first** — builds on macOS, Linux, and BSD; Windows is not supported because file locking uses `syscall.Flock`

## Tech Stack

| Layer | Technology | Purpose |
| --- | --- | --- |
| Language | [Go](https://go.dev) | Single static binary, no runtime dependencies |
| CLI framework | [Cobra](https://github.com/spf13/cobra) | Command/subcommand tree, flag parsing, help generation |
| Flag parsing | [pflag](https://github.com/spf13/pflag) | POSIX/GNU-style flags (`-a`, `--all`) |
| Persistence | `encoding/csv` (stdlib) | Read/write the `todo.csv` data file |
| File locking | `syscall.Flock` (stdlib) | Advisory lock around read-modify-write cycles (Unix only) |
| Output formatting | `text/tabwriter` (stdlib) | Aligned column rendering in `list` |
| Time formatting | [timediff](https://github.com/mergestat/timediff) | Relative timestamps ("a few seconds ago") |

## Requirements

- Go **1.26** or newer (see `go.mod`)
- macOS, Linux, or a BSD — Windows builds fail today because `syscall.Flock` is Unix-only

## Installation

### From source (recommended for local use)

```sh
git clone <repository-url> todo-cli
cd todo-cli
go build -o task ./cmd/task
```

The binary is produced as `./task` in the project root.

### Install into `$GOBIN`

```sh
go install ./cmd/task
```

This places `task` in `$(go env GOPATH)/bin`, which is normally already on your `PATH`.

> **Note:** the entrypoint lives in `cmd/task`, so build with `./cmd/task` — a bare `go build .` will not work from the repository root.

### Verify the installation

```sh
task --help
```

## Usage

```
task [command]
```

### Commands

| Command | Description | Example |
| --- | --- | --- |
| `task add <task>` | Add a new task | `task add "write release notes"` |
| `task list` | List pending tasks | `task list` |
| `task list --all` | List all tasks, including completed | `task list -a` |
| `task complete <id>` | Mark a task as done | `task complete 3` |
| `task delete <id>` | Delete a task permanently | `task delete 3` |
| `task help` | Show help for any command | `task help list` |

### Workflow example

```sh
$ task add "buy coffee beans"
$ task add "review pull request #42"
$ task list
ID          Task                             Created                    Done
1           buy coffee beans                 a few seconds ago          false
2           review pull request #42          a few seconds ago          false

$ task complete 1
$ task list --all
ID          Task                             Created                    Done
1           buy coffee beans                 a few seconds ago          true
2           review pull request #42          a few seconds ago          false

$ task delete 2
$ task list
ID          Task          Created          Done
```

`add`, `complete`, and `delete` are silent on success; only `list` prints output.

## Data Storage

Tasks are persisted to **`todo.csv` in the current working directory**, so the file is created next to wherever you run the command.

> Run `task` from the same directory each time, or the CLI will look for a different (or empty) data file.

File format (values below are illustrative; `add` always writes `Done=false`):

```csv
No,Task,CreatedAt,Done
1,buy coffee beans,2026-10-28 18:48:54,false
2,"review pull request #42",2026-10-28 19:47:43,true
```

| Column | Type | Description |
| --- | --- | --- |
| `No` | integer | Task identifier, used by `complete` and `delete` |
| `Task` | string | Task description (may contain commas or newlines; quoted per RFC 4180) |
| `CreatedAt` | timestamp | Creation time in `Asia/Jakarta` (WIB), format `YYYY-MM-DD HH:MM:SS` |
| `Done` | boolean | Completion state (`true` / `false`) |

`todo.csv` is listed in `.gitignore` and is intentionally **not** committed — it is local user data, not source code.

## Project Structure

```
todo-cli/
├── cmd/
│   ├── root.go        # root command, flag and subcommand wiring
│   ├── add.go         # task add
│   ├── list.go        # task list (--all flag)
│   ├── complete.go    # task complete
│   ├── delete.go      # task delete
│   └── task/
│       └── main.go    # entrypoint (package main)
├── internal/
│   └── todo.go        # domain model, CSV persistence, file locking, rendering
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

The `cmd` package contains only command definitions and flag handling; all domain logic and I/O live in `internal`.

## Development

```sh
go build ./...        # compile all packages
go vet ./...          # static analysis
gofmt -l cmd internal # list unformatted files
```

## Known Limitations

These are current behaviours, verified against the built binary:

- **`delete` can panic** on a non-last element — deleting any task other than the last one raises `index out of range` because the slice is mutated while being iterated. Deleting the final task works.
- **Invalid or unknown IDs fail silently** — `task complete abc` prints the raw `strconv.Atoi` error twice and `task complete 99` prints nothing; both exit `0` and leave the file untouched.
- **The lock does not cover the write** — `saveTodo` reopens the file with `os.Create`, so the lock acquired on the read handle protects readers only, not the rewrite.
- **Empty tasks are accepted** — `task add ""` writes a row with an empty description.
- **IDs may be reused** — new IDs are derived from `len(tasks) + 1`, so removing a task can cause the next `add` to reuse its number.
- **Writes are not atomic** — `todo.csv` is reopened with `os.Create` and rewritten in place; an interruption mid-write can truncate the file.
- **Data path is relative to the working directory**, not resolved per user; see [Data Storage](#data-storage).
- **Help text is partly Indonesian** (`add`, `list`) and partly English (`complete`, `delete`).
- **`complete` and `delete` reject missing arguments** via `cobra.ExactArgs(1)`, so the ID is mandatory despite `Usage:` lines not showing it.

## License

See [LICENSE](LICENSE).
