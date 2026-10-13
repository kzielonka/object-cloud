# Go Cheat Sheet & Mental Model (TypeScript/Node to Go)

This cheat sheet maps concepts from TypeScript/Jest (like `__repoSharedTests__/shopifyAccountRepoSharedTests.ts`) to idiomatic Go.

---

## 1. Test Discovery & File Naming Rules

In Jest/TypeScript, test runners match glob patterns like `**/*.test.ts`, `**/*.spec.ts`, or anything inside `__tests__/` or `__repoSharedTests__/`. In Go, the compiler enforces strict rules:

| Concept | TypeScript / Jest | Go |
| :--- | :--- | :--- |
| **Test file suffix** | `.test.ts` or `.spec.ts` | **Must be `_test.go` (singular)** (e.g. `storage_test.go`) |
| **Plural `_tests.go`** | Treated as test file if matching glob | ❌ **Treated as production code!** (Breaks build if package names mismatch) |
| **Test function signature** | `it("does something", () => {})` | `func TestXxx(t *testing.T)` (Starts with `Test` + Capital letter) |
| **Shared / Helper functions** | `describeSharedTests(...)` | `func RunSharedContractTests(t *testing.T, ...)` (Does not start with `Test` so `go test` won't run it directly) |

---

## 2. Package Organization & Directory Rules

In TypeScript, multiple files in the same folder can export anything independently. In Go:

1. **One package per folder:** Every non-test file in a directory **must** declare the same package name (e.g., `package object`).
2. **Never name a file `_tests.go`:** If you put `package object_test` in a file named `storage_contract_tests.go`, Go treats it as production code and complains:
   ```text
   found packages object (storage.go) and object_test (storage_contract_tests.go) in ...
   ```
3. **The `_test` package exception (Black-Box vs White-Box Testing):**
   Within any test file (`*_test.go`), Go allows **only two choices** for the package name:

| Testing Style | Package Declaration | What it can access | When to use |
| :--- | :--- | :--- | :--- |
| **White-Box** (Internal) | `package auth` (same as directory) | **Everything:** both exported (`Public`) and unexported (`private`) structs, fields, functions | Testing complex unexported logic, mocking internal state, testing private helpers |
| **Black-Box** (External) | `package auth_test` (`<pkg>_test`) | **Public API only:** must import the package (`import ".../internal/auth"`), can only see exported symbols | Testing the package from a real consumer's viewpoint; **prevents circular dependencies** |

---

## 3. Visibility & Type Syntax (Public vs Private)

In Go, there are no `public` or `private` keywords. Visibility is strictly controlled by **capitalization**:

| Concept | Unexported (Private) | Exported (Public) |
| :--- | :--- | :--- |
| **Identifiers** | Lowercase first letter: `uploadPolicy`, `key` | Uppercase first letter: `UploadPolicy`, `Key` |
| **Scope** | Visible only inside the same package | Accessible to any importing package |
| **JSON Serialization** | `json.Marshal` ignores lowercase fields! | Must be uppercase to be serialized to JSON |

### Struct Declaration & Instantiation Syntax:
* **Declare a struct:** use `type Name struct`, **never** `struct Name` (which is a syntax error).
* **Return pointer:** Instantiate and return address directly without temporary variables:

```go
// 1. Declaration
type UploadPolicy struct {
    Key string `json:"key"`
}

// 2. Constructor returning pointer
func NewUploadPolicy(key string) *UploadPolicy {
    return &UploadPolicy{
        Key: key,
    }
}
```

### Struct Tags (Go's Metadata Annotations):
Struct tags are compile-time string metadata in backticks (`` `key:"value"` ``) attached to struct fields.

* **History & Philosophy:** Present since Go's 2009 release (and Go 1.0). The creators (Rob Pike, Ken Thompson) deliberately avoided heavy annotation/macro frameworks (like in Java/Spring). They chose the simplest possible compromise: raw string metadata attached to fields, inspected at runtime via reflection (`reflect.StructTag`).
* **Not Just for JSON:** It is a general-purpose language feature:
  * `json:"key"` &rarr; `encoding/json`
  * `db:"col_name"` / `gorm:"primaryKey"` &rarr; Database ORMs
  * `validate:"required,email"` &rarr; Struct validation libraries
  * `env:"PORT"` &rarr; Environment variable loaders
* **Combining tags:** Multiple packages can read the same field:
  ```go
  type UploadPolicy struct {
      Key       string `json:"key" db:"object_key" validate:"required"`
      ExpiresAt int64  `json:"expires_at,omitempty"`
      Secret    string `json:"-"` // "-" means completely ignore / hide from JSON
  }
  ```
* **Golden Rule for JSON:** The field **must start with a Capital letter (exported)**, or `json.Marshal` will silently ignore it regardless of any tags!

### Functional Options Pattern (The Go Alternative to OOP Builders):
Instead of Java-style fluent builders (`new Builder().WithX().Build()`), Go uses variadic closures (`opts ...Option`) for constructors with optional configuration:

```go
type UploadOption func(*UploadPolicy)

func WithMaxContentLength(bytes int64) UploadOption {
    return func(p *UploadPolicy) {
        p.MaxContentLength = bytes
    }
}

// Usage: NewUploadPolicy("avatar.png", WithMaxContentLength(5*1024*1024))
func NewUploadPolicy(key string, opts ...UploadOption) *UploadPolicy {
    p := &UploadPolicy{Key: key}
    for _, opt := range opts {
        opt(p)
    }
    return p
}
```

---

## 4. Reusable Contract / Compliance Tests Pattern

### TypeScript / Jest Pattern:
```typescript
// __repoSharedTests__/shopifyAccountRepoSharedTests.ts
export function describeShopifyAccountRepoContract(createRepo: () => ShopifyAccountRepo) {
  describe("ShopifyAccountRepo Contract", () => {
    it("saves and finds account", async () => {
      const repo = createRepo();
      // assertions...
    });
  });
}
```

### Equivalent Idiomatic Go Pattern:
```go
// internal/object/storage_contract_test.go  (Notice: _test.go singular)
package object_test

import (
	"testing"
	"github.com/kzielonka/object-cloud/internal/object"
)

type FileSystemFactory func(t *testing.T) object.FileSystem

// RunFileSystemContract is a helper (does NOT start with Test), so it only runs when invoked
func RunFileSystemContract(t *testing.T, newFS FileSystemFactory) {
	t.Run("saves and reads back file content", func(t *testing.T) {
		fs := newFS(t)
		// test logic...
	})

	t.Run("returns ErrNotFound when file does not exist", func(t *testing.T) {
		fs := newFS(t)
		// test logic...
	})
}
```

### Running the Contract in Adapter Test Files:
```go
// internal/object/in_memory_storage_test.go
func TestInMemoryFileSystem_Contract(t *testing.T) {
	RunFileSystemContract(t, func(t *testing.T) object.FileSystem {
		return object.NewInMemoryFileSystem()
	})
}
```

---

## 5. Errors & Comparisons

| Action | TypeScript | Go |
| :--- | :--- | :--- |
| **Throwing error** | `throw new Error("msg")` | `return nil, errors.New("msg")` |
| **Wrapping error** | `new Error("msg", { cause: err })` | `fmt.Errorf("msg: %w", err)` |
| **Catching / Checking** | `if (err instanceof NotFoundError)` | `if errors.Is(err, object.ErrNotFound)` |
| **Extracting error struct** | `const e = err as CustomError` | `var e *CustomError; if errors.As(err, &e) { ... }` |

---

## 6. Quick Syntax Differences

* **No semicolons:** Do not put `;` at the end of lines.
* **No colons in parameter lists:** `func(path string)` (NOT `path: string`).
* **Struct initialization:** `&MyStruct{Field: val}` with curly braces `{}` (NOT `new MyStruct()` or `MyStruct()`).
* **Byte comparison:** `bytes.Equal(a, b)` instead of converting `string(a) == string(b)`.
