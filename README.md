FluentFGA
=========

Code generation tool for OpenFGA authorization models.

# Install

```shell
go get github.com/sectrean/fluentfga
```

# Usage

Configure code generation to generate types from your OpenFGA authorization model.

```go
//go:generate go run github.com/sectrean/fluentfga/cmd/fluentfga generate --clean model.fga ./model/
```

# Example

```go
var client sdkclient.SdkClient

user := model.User{ID: "anne"}
commenter := model.DocumentCommenterRelation{}
doc := model.Document{ID: "2021-budget"}

// Write
err := fluentfga.Write(
    commenter.NewTuple(user, doc),
).Execute(ctx, client)

// Check
allowed, err := fluentfga.Check(
    user,
    commenter,
    doc,
).Execute(ctx, client)

// ListObjects
var docs []model.Document
docs, err := fluentfga.ListObjects(
    user,
    commenter,
).Execute(ctx, client)

// ListUsers
var users []model.User
users, err := fluentfga.ListUsers(
    doc,
    commenter,
    fluentfga.UserTypeFilter[model.User]{},
).Execute(ctx, client)
```

# TODO

- Implement BulkCheck operation.
- Implement support for Conditions.
- Allow customization of generated code:
    - Type names
    - ID field names
    - ID field types
    - Package name
