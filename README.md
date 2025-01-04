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
//go:generate go run github.com/sectrean/fluentfga/cmd/fluentfga generate --clean model.fga ./
```

# Examples

```go
var client sdkclient.SdkClient

beth := model.User{ID: "beth"}
commenter := model.DocumentCommenterRelation{}
doc := model.Document{ID: "2021-budget"}

// Write
err := fluentfga.Write(
    commenter.NewTuple(beth, doc),
).Execute(ctx, client)

// Check
allowed, err := fluentfga.Check(
    beth,
    commenter,
    doc,
).Execute(ctx, client)

// ListObjects
var docs []model.Document
docs, err := fluentfga.ListObjects(
    beth,
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

- Support `fga.mod` and `.json` model files.
- Implement BulkCheck operation.
- Implement support for Conditions.
- Allow customization of generated code:
    - Type names
    - ID field names
    - ID field types
    - Package name
