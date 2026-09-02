# Next Steps & Technical Roadmap

This document tracks upcoming architectural tasks, optimizations, and technical debt for `object-cloud`.

---

## 1. Key Hashing & Storage Path Sanitization (`Store`)
- [ ] **v1: Flat Key Hashing in `Store`:**
  - *Context:* Raw user keys (e.g., `../../etc/passwd`, `photos/2026:01*?.png`, long URLs) can cause path traversal vulnerabilities and filesystem errors.
  - *Action:* `Store` will hash user keys with SHA-256 (via `crypto/sha256` and `encoding/hex`) before passing the safe, deterministic 64-character hex string to `FileSystem`.
  - *Security & Compatibility:* Guarantees flat, safe alphanumeric filenames on any OS filesystem and prevents directory traversal attacks.
- [ ] **v2 (Future): Prefix Sharding / Fan-out:**
  - *Context:* Storing hundreds of thousands of files in a single flat directory causes filesystem inode lookup degradation.
  - *Action:* Split hash into subdirectories (e.g. `hash[:2] / hash[2:]` like Git objects).

---

## 2. Safety & Error Handling in `discFileSystem.SaveFile`
- [ ] **Check `os.Create` error before `defer outFile.Close()`:**
  - *Context:* If `os.Create(fullPath)` fails (e.g. permission denied, disk full), `outFile` is `nil`. Calling `defer outFile.Close()` panics with a nil pointer dereference.
  - *Action:* Check `if err != nil { return err }` before deferring `outFile.Close()`.
- [ ] **Handle `io.Copy` errors:**
  - *Context:* Streaming write may fail mid-transfer if connection drops or disk runs out of space.
  - *Action:* Return and wrap any error returned by `io.Copy(outFile, data)`.

---

## 3. True Streaming Downloads (`io.ReadCloser`)
- [ ] **Avoid buffering entire files into memory in `OpenFile`:**
  - *Context:* Currently, `discFileSystem.OpenFile` reads the entire file into RAM with `io.ReadAll(file)` and returns `bytes.NewReader(data)`. For large files (e.g., 100MB - 5GB), this consumes high memory.
  - *Action:* Upgrade `FileSystem.OpenFile` interface to return `(io.ReadCloser, error)` so consumers can stream directly from `*os.File` and close it when finished:
    ```go
    type FileSystem interface {
        SaveFile(path string, data io.Reader) error
        OpenFile(path string) (io.ReadCloser, error)
    }
    ```
  - *Adapter updates:*
    - `discFileSystem.OpenFile`: Return `file, nil` directly (`*os.File` implements `io.ReadCloser`).
    - `inMemoryFileSystem.OpenFile`: Wrap `bytes.NewReader(data)` with `io.NopCloser(bytes.NewReader(data))`.
