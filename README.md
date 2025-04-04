# DynamicQueryKit (DQK) 🧠⚡️

A modular, lightweight Go library to help you build powerful, flexible APIs with dynamic SQL filtering, clean decoding, permission logic, and plug-and-play middlewares.

---

## ✨ Features

- ✅ Decode request bodies (JSON & XML) with detailed error handling
- ✅ Build dynamic SQL filters using [Squirrel](https://github.com/Masterminds/squirrel)
- ✅ Middleware: CORS, Logging, Middleware stacking
- ✅ Permission-based access control with scoped and global checks
- ✅ Pagination helper that wraps Squirrel queries
- ✅ Idiomatic Go, clean API, testable and composable components

---

## 📦 Installation

```bash
go get github.com/yourusername/dynamicquerykit
```

---

## 🔍 Quick Examples

### 🔸 DecodeBody

```go
type User struct {
	ID   int    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
}

var user User
status, err := dynamicquerykit.DecodeBody("application/json", r.Body, &user)
```

Supports JSON & XML decoding with descriptive HTTP error codes and messages.

---

### 🔸 DataEncode

Respond dynamically based on `Accept` header (JSON, XML).

```go
data := map[string]string{"message": "hello"}
bytes, contentType, _ := dynamicquerykit.DataEncode("application/json", data)
w.Header().Set("Content-Type", contentType)
w.Write(bytes)
```

---

### 🔸 Permission Utilities

```go
perms := dynamicquerykit.BuildPermissionSet([]string{"view_users", "edit_users"})

if dynamicquerykit.HasPermission("edit_users", perms) {
	// allow edit
}

if dynamicquerykit.IsAccessRestricted("view_users", "view_all_users", perms) {
	http.Error(w, "forbidden", http.StatusForbidden)
}
```

---

### 🔸 Pagination Helper (Squirrel)

```go
query := squirrel.Select("id").From("cars").Where(squirrel.Eq{"color": "black"})
total, err := dynamicquerykit.GetPaginationData(ctx, db, &query)
```

Wraps the query in a `SELECT COUNT(*) FROM (...)` for total count pagination support.

---

### 🔸 Middleware Stack

```go
stack := dynamicquerykit.CreateStack(
	dynamicquerykit.Cors,
	dynamicquerykit.Logging,
)

http.Handle("/", stack(http.HandlerFunc(myHandler)))
```

Use `CreateStack` to compose middlewares easily, including built-in `Logging` and `Cors`.

---

## 🥪 Testing

Tests cover all utilities:

- Body decoding (valid/invalid/malformed input)
- Permission logic
- Middleware execution
- Pagination query generation

You can run all tests with:

```bash
go test ./...
```

---

## 🧠 Philosophy

DynamicQueryKit is designed to be:

- **Minimal**: No unnecessary dependencies
- **Composable**: Everything is made to work together or independently
- **Practical**: Targets real problems when building APIs in Go

---
