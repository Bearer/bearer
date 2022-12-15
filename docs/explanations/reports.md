---
title: Report Types
layout: layouts/doc.njk
---

# Report Types

Curio can generate four types of reports about your codebase, all from the same underlying scan.

## Policy Report

Policies are the core feature of Curio. The policy report allows you to quickly see data security vulnerabilities in your codebase. The report breaks down policy violations by severity level: Critical, High, Medium, Low. These levels help you prioritize issues and fix the most important vulnerabilities.

For each violation, the policy report includes the affected file and, when possible, the line of code and a snippet of the surrounding code. Here's an excerpt from the policy report run on our [example publishing app](https://github.com/Bearer/bear-publishing):

```txt
...
CRITICAL: JWT leaking policy failure with Personal data
JWT leaks detected. Avoid storing sensitive data in JWTs.

File: lib/jwt.rb:6

 3 JWT.encode(
 4       {
 5         id: user.id,
 **6         email: user.email,**
 7         class: user.class,
 8       },
 9       nil,
 10       'HS256'
 11     )

...

=====================================

Policy failures detected

14 policies were run and 6 failures were detected.

CRITICAL: 0
HIGH: 4 (Insecure HTTP with Data Category, Logger leaking, JWT leaking, Cookie leaking)
MEDIUM: 2 (Insecure SMTP, Insecure FTP)
LOW: 0
```

Policy reports are [currently available](/reference/supported-languages/) for Ruby projects, with Javascript/Typescript coming soon and more languages to follow.

To run your first policy report, run `curio scan` with the `--report policies` flag.

## Data Flow Report

The data flow report breaks down the data types and associated components detected in your code. It focuses your app processes personal and sensitive data and where it may be exposed to third parties and databases.

This report type can produce JSON or YAML (using the `--format` flag) and is best used for identifying where data exists in your code. You can then use this to create a record of processing activity (ROPA), AppStore report, or build your data catalog. In the following example, we can see all the places an `Email Address` is processed by our [example application](https://github.com/Bearer/bear-publishing):

```bash
{
      "name": "Email Address",
      "detectors": [
        {
          "name": "rails",
          "locations": [
            {
              "filename": "db/schema.rb",
              "line_number": 91
            }
          ]
        },
        {
          "name": "ruby",
          "locations": [
            {
              "filename": "app/controllers/application_controller.rb",
              "line_number": 35
            },
            {
              "filename": "app/controllers/application_controller.rb",
              "line_number": 37
            },
            {
              "filename": "app/controllers/application_controller.rb",
              "line_number": 42
            },
            {
              "filename": "app/controllers/orders_controller.rb",
              "line_number": 105
            },
            {
              "filename": "app/models/user.rb",
              "line_number": 8
            },
            {
              "filename": "app/services/marketing_export.rb",
              "line_number": 23
            },
            {
              "filename": "lib/jwt.rb",
              "line_number": 6
            }
          ]
        }
      ]
    }
```

If we look at the `db/schema.rb` file mentioned in the report, we can see that email is exposed:
```ruby
  create_table "users", force: :cascade do |t|
    t.string "name"
    t.string "email"
    t.string "telephone"
    t.integer "organization_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["organization_id"], name: "index_users_on_organization_id"
  end
```

To run your first data flow report, run `curio scan` with the `--report dataflow` flag.

## Stats Report

The stats report provides a minimal overview of the scan including the total lines scanned, the detected data types, and the number of times each was found.

```bash
Scanning target... 100% [========================================] (278/278, 103 files/s) [2s]
{
  "number_of_lines": 1456,
  "number_of_data_types": 6,
  "data_types": [
    {
      "name": "Email Address",
      "occurrences": 8
    },
    {
      "name": "Fullname",
      "occurrences": 3
    },
    {
      "name": "Physical Address",
      "occurrences": 2
    },
    {
      "name": "Telephone Number",
      "occurrences": 1
    },
    {
      "name": "Transactions",
      "occurrences": 7
    },
    {
      "name": "Unique Identifier",
      "occurrences": 2
    }
  ]
}
```

This report is useful for putting together ongoing statistics, making a quick check of data processing, or building internal dashboards. For more detailed information, you'll want to run a data flow report.

To run your first policy report, run `curio scan` with the `--report stats` flag.

## Detectors Report

The detectors report type is the most low-level, data-rich type. You’re unlikely to use this report on its own, but it can be useful for building your own tooling based on the data parsed by Curio.

To run your first detectors report, run `curio scan` with the `--report detectors` flag.
