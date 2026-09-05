# Next Steps & Architectural Roadmap

This document outlines the architectural decisions, technical debt, and future roadmap for `object-cloud`.

---

## 1. Project & Repository Strategy

* **Modular Monorepo (Single Repo, Pluggable Packages):**
  * Rather than splitting into multiple micro-repositories, maintain a modular Go layout.
  * `pkg/object`: Reusable public Go library (exportable for other Go projects via `go get`).
  * `internal/server`: HTTP API handlers and routing.
  * `internal/image` (future): Optional image resizing / ImageMagick processing layer.
  * `cmd/server/main.go`: Server daemon entrypoint.

---

## 2. Storage Engine (Single-Node Primitive)

The single-node storage engine (`pkg/object`) serves as the foundational drive-level primitive (similar to MinIO's `xl-storage` or CockroachDB's `Pebble`), managing raw I/O on attached disk storage.

### Key Hashing & Sanitization (`Store`)
- [ ] **v1 Flat Key Hashing:**
  * Logical user keys (e.g., `avatars/user/photo.jpg`, `../../etc/passwd`) are hashed with SHA-256 (`crypto/sha256` + `encoding/hex`) in `Store`.
  * Guarantees 100% path-traversal protection and uniform safe filenames across all operating systems.
- [ ] **v2 Prefix Sharding / Fan-out:**
  * When scaling to millions of files, split hash into subdirectories (e.g. `/data/e3/b0c442...` like Git objects) to prevent single-directory inode performance bottlenecks.

### Crash Safety & Atomic Writes (`diskFileSystem`)
- [ ] **Check `os.Create` error before `defer outFile.Close()`:**
  * Prevents nil pointer dereference panics when file creation fails.
- [ ] **Atomic Writes via Temp File & `os.Rename`:**
  * Write incoming stream to a temporary file (`.tmp_...`) first, then atomically rename to target path upon successful completion to prevent corrupt/partial files on disk if the connection or server drops.
- [ ] **Handle `io.Copy` errors:**
  * Ensure streaming write errors are properly caught and bubbled up.

### True Streaming Downloads (`io.ReadCloser`)
- [ ] **Avoid buffering in RAM:**
  * Upgrade `FileSystem.OpenFile` interface to return `(io.ReadCloser, error)`.
  * Stream `*os.File` directly on disk and use `io.NopCloser` for in-memory byte readers.

---

## 3. Data Integrity: HTTP ETags vs Bitrot Protection

* **Bitrot Protection:** Deliberately **skipped**. Modern SSDs, cloud storage (AWS EBS, GCP PD), and modern filesystems already have hardware-level ECC. Software-level bitrot detection adds CPU overhead without value on single-node setups.
* **HTTP `ETag` Headers:** Compute and return standard SHA-256 / MD5 checksums in `ETag` response headers so HTTP clients/browsers can verify transfer integrity.

---

## 4. HTTP API & Server Layer (`internal/server`)

- [ ] **Upload Endpoint (`PUT /objects/{key...}`):**
  * Stream request body directly into `store.Upload(key, r.Body)`.
  * Return `201 Created` with `ETag`.
- [ ] **Download Endpoint (`GET /objects/{key...}`):**
  * Stream stored object to `w` via `io.Copy(w, stream)`.
  * Return `404 Not Found` when `errors.Is(err, object.ErrNotFound)`.
- [ ] **HTTP Testing:**
  * Use Go's built-in `net/http/httptest` package (`httptest.NewServer` and `httptest.ResponseRecorder`) for full endpoint testing.

---

## 5. Future: Distributed Cluster Layer (`pkg/cluster`)

The single-node `pkg/object` can later be wrapped by a cluster coordination layer:
* **Node Discovery / Gossip:** Tracking cluster membership.
* **Consistent Hashing Ring:** Mapping keys to specific storage nodes.
* **Quorum Replication:** Writing concurrently to multiple nodes (e.g. write to 3 nodes, require 2 acks).
* **Self-Healing:** Re-syncing missing or out-of-date files when a node recovers.
