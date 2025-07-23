# gocrate

**gocrate** is a tool for packaging, distributing, and managing self-contained binary crates in JSON format. It supports creating, installing, and uninstalling binary distributions with embedded metadata and binary content. The project is written in Go and targets Linux systems.

> **License**: GNU General Public License v3.0 (GPLv3)

---

## Overview

A `crate` is a single JSON file containing:

- Binary file contents (embedded)
- Metadata: project name, binary filename
- Optional source URL

The tool supports:

- Crate building from a local binary
- Installation to system-wide paths (`/bin`)
- Uninstallation from system paths
- Binary extraction to a user-defined prefix
- Pulling and processing crates from remote URLs

---

## Features

- ✅ Easy creates building from local binaries
- ✅ All data saved as a single JSON file
- ✅ System installation and uninstallation of binaries
- ✅ Remote crate fetching
- ✅ Interactive command-line interface
- ❌ No support for colorized output (yet)
- ❌ No support for Windows installation (yet)

---

### Build a crate

```bash
$ ./gocrate build
Enter the project name: mytool
Enter the binary file name: ./mytool
Enter the source URL (optional, press Enter to skip): https://example.com/mytool
````

Produces `mytool.json`.

---

### Install a crate (requires root)

```bash
$ sudo ./gocrate install mytool.json
```

Installs `mytool` to `/bin/mytool`.

---

### Uninstall

```bash
$ sudo ./gocrate uninstall mytool.json
```

Removes `/bin/mytool`.

---

### Unpack to a custom path

```bash
$ ./gocrate get-bin mytool.json
Enter the prefix path to unpack the binary: ./bin/
```

Places the binary at `./bin/mytool`.

---

### Pull from remote URL

```bash
$ sudo ./gocrate pull https://example.com/mytool.json
Do you want to install or uninstall the crate? (install/uninstall): install
```

---

## Implementation Notes

* Uses Go standard libraries: `os`, `io`, `encoding/json`, `net/http`, `os/exec`
* Crate files are portable JSON blobs containing full binary content
* Installation/uninstallation directly modifies `/bin`, not PATH-managed locations
* Interactive command interface with minimal validation

---

## System Requirements

* Go 1.18+
* Unix-like OS (preferably Linux)
* Root access (via `sudo`) for system-level operations

---

## License

This project is licensed under the **GNU General Public License v3.0** (GPLv3).
See [`LICENSE`](./LICENSE) for details.

---

## Roadmap

* [x] Local crate building
* [x] Binary install/uninstall
* [x] HTTP crate fetching
* [ ] Colorful output 
* [ ] Better command system
* [ ] Virus scanning (via virus total)
* [ ] Windows support

---