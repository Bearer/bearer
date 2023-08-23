---
title: Add or update a recipe
---

# Add or update a recipe

Recipes are part of how Bearer CLI makes connections between your code and other sources. These are things like data stores, APIs, and internal services. They work by providing information about endpoints, API base urls, package information, etc.

Recipes are located at `bearer/internal/classification/db/recipes/`.

```md
.
├ internal/
│  └ classification/
│       └  db/
│           └  recipes/
```

Each recipe consists of a `JSON` file containing the following properties:

- `metadata` (object): Metadata about the recipe itself. Used for tracking recipe versions. Ex: `"metadata": { "version": "1.0" }`.
- `name` (string): The name of the recipe.
- `type` (string): The recipe type. Supported  types are `external_service`, `internal_service`, and `data_store`.
- `urls` (array of strings): URLs used by the service. Bearer CLI will use this to aid in finding connections within a codebase. Supports wildcards.
- `exclude_urls` (array of strings): Any urls that would be caught by the wildcards in the `urls` list that you'd like to explicitly exclude.
- `packages` (array of objects): Common packages that connect to the service. Each package object should contain:
  - `name` (string): The official name of the package used by package managers.
  - `group` (string): For Java applications (e.g., `maven`). Set to `null` for other use cases.
  - `package_manager` (string): The package manager that manages the package, such as npm, go, etc.
- `uuid`: A unique identifier to distinguish the recipe from others. See below for [generating a new uuid](#generating-a-uuid).
- `sub_type` (string): The subtype of the earlier `type` property.
  - `external_service` subtypes:
    - `third_party`
  - `data_store` subtypes:
    - `database`
    - `datalake`
    - `object_storage`
    - `search_engine`
    - `key_value_cache`
    - `flat_file`
    - `shared_folders`
  - `internal_service` subtypes:
    - `message_bus`

If any of the existing properties and available values don't meet the needs of your new recipe, [open a new issue]({{meta.sourcePath}}/issues/new/choose). You can view the existing recipes [in the GitHub repo]({{meta.sourcePath}}/tree/main/internal/classification/db/recipes).

## Generating a UUID

Recipes each contain a universally unique identifier (UUID). To generate one, use the `uuidgen` tool.

On linux:

```bash
uuidgen
```

On MacOS, you need to force the output to lowercase:

```bash
uuidgen | tr "[:upper:]" "[:lower:]"
```

## Commiting the new recipe

To contribute the new recipe to Bearer CLI, refer to the [Contributing Code guide](/contributing/code/).
