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
2. **The `_test` package exception:** Files ending in `_test.go` can declare `package object_test` (for external black-box testing) or `package object` (for internal white-box testing).
3. **Never name a file `_tests.go`:** If you put `package object_test` in a file named `storage_contract_tests.go`, Go treats it as production code and complains:
   ```text
   found packages object (storage.go) and object_test (storage_contract_tests.go) in ...
   ```

---

## 3. Reusable Contract / Compliance Tests Pattern

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

// RunStorageContractTests is a helper (does NOT start with Test), so it only runs when invoked
func RunStorageContractTests(t *testing.T, newFS FileSystemFactory) {
	t.Run("saves and opens file", func(t *testing.T) {
		fs := newFS(t)
		// test logic...
	})

	t.Run("returns ErrNotFound when missing", func(t *testing.T) {
		fs := newFS(t)
		// test logic...
	})
}
```

### Running the Contract in Adapter Test Files:
```go
// internal/object/in_memory_storage_test.go
func TestInMemoryFileSystem_Contract(t *testing.T) {
	RunStorageContractTests(t, func(t *testing.T) object.FileSystem {
		return object.InMemoryFileSystem()
	})
}
```

---

## 4. Errors & Comparisons

| Action | TypeScript | Go |
| :--- | :--- | :--- |
| **Throwing error** | `throw new Error("msg")` | `return nil, errors.New("msg")` |
| **Wrapping error** | `new Error("msg", { cause: err })` | `fmt.Errorf("msg: %w", err)` |
| **Catching / Checking** | `if (err instanceof NotFoundError)` | `if errors.Is(err, object.ErrNotFound)` |
| **Extracting error struct** | `const e = err as CustomError` | `var e *CustomError; if errors.As(err, &e) { ... }` |

---

## 5. Quick Syntax Differences

* **No semicolons:** Do not put `;` at the end of lines.
* **No colons in parameter lists:** `func(path string)` (NOT `path: string`).
* **Struct initialization:** `&MyStruct{Field: val}` with curly braces `{}` (NOT `new MyStruct()` or `MyStruct()`).
* **Byte comparison:** `bytes.Equal(a, b)` instead of converting `string(a) == string(b)`.
