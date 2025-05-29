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
or even better. Add the handlers to your server mux

```go
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", controllers.GetHealth)

	middlewareStack := dqk.CreateStack(dqk.Logging, dqk.Cors)

	server := http.Server{
        Addr:    ":8080",
		Handler: middlewareStack(mux),
	}

	server.ListenAndServe()
}

```

Use `CreateStack` to compose middlewares easily, including built-in `Logging` and `Cors`.

---

## 🧠 Philosophy

DynamicQueryKit is designed to be:

- **Minimal**: No unnecessary dependencies
- **Composable**: Everything is made to work together or independently
- **Practical**: Targets real problems when building APIs in Go. Especially dynamic queries

---


## Example workflow



- Specify filters

```go
filters := []dqk.filers{
    {Name: "color", Operator: "IN", DbField: "colors.name", FieldID: "colors.id"},                                                      // added FieldID
    // ... any other columns
}
```

- Make your base query dqk uses squirrel for query building
```go
myquery := sq.Select("id,name").From("cloths as c").Join("colors as cl ON cl.id = c.color_id")
```

- Get user **request parameters** -> params["color"] = ["blue","red"]
- Pass the params and the base query with your filters to dqk.DynamicFilters

```go
query, _ := dqk.DynamicFilters(filters, myquery, r.URL.Query()) //or pass your params variable

// query = sq.Select("id,name").From("cloths as c").Join("colors as cl ON cl.id = c.color_id").Where("color IN (blue,red)")
```

- (optional) if you want pagination 
```go
paginationQuery := dqk.GetPaginationQuery(query)
// do not limit the pagination query, you will only get the number you specified as limit
query = query.limit(mylimit).offset(myoffset)
query = query.OrderBy(dqk.OrderValidation(r.URL.Query().get("order_by"),r.URL.Query().Get("order_direction"),filters))
// Optionally you can specify nulls last
// query = query.OrderBy(dqk.OrderValidation(r.URL.Query().get("order_by"),r.URL.Query().Get("order_direction"),filters + " NULLS LAST"))
```





