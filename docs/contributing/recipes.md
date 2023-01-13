---
title: Add or update a recipe
---

# Add or update a recipe

Recipes are part of how Curio makes connections between your code and other sources. These are things like data stores, APIs, and internal services. They work by providing information about endpoints, API base urls, package information, etc.

Recipes are located at `curio/pkg/classification/db/recipes/`

```
.
├—— pkg/
│    └—— classification/
│         └——  db/
│               └——  recipes/
```

Each recipe consists of a `JSON` file containing the following properties:

- `metadata` (object): Metadata about the recipe itself. Used for tracking recipe versions. Ex: `"metadata": { "version": "1.0" }`.
- `name` (string): The name of the recipe.
- `type` (string): The recipe type. Supported  types are `external_service`, `internal_service`, and `data_store`.
- `urls` (array of strings): URLs used by the service. Curio will use this to aid in finding connections within a codebase. Supports wildcards.
- `exclude_urls` (array of strings): Any urls that would be caught by the wildcards in the `urls` list that you'd like to explicitly exclude.
- `packages` (array of objects): Common packages that connect to the service. Each package object should contain:
  - `name`: The official name of the package used by package manager files to identify it.
  - `group`: null
  - `package_manager`: The package manager that manages the package, such as npm, go, etc.
- `uuid`: A unique identifier to distinguish the recipe from others.
- `sub_type`: The sub type of the earlier `type` property.

