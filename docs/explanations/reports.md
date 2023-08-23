---
title: Report Types
---

# Report Types

Bearer CLI can generate various types of reports about your codebase, all from the same underlying scan.

## Security Report

- Usage: `bearer scan . --report security`
- Default format: `json`

The security report allows you to quickly see security risks and vulnerabilities found in your codebase using a security [scanner type](/explanations/scanners) (SAST by default).

For each violation, the report includes the affected file and, when possible, the line of code and a snippet of the surrounding code. Here's an excerpt from the security report run on the [OWASP Juice Shop app](https://github.com/juice-shop/juice-shop):

```txt
$ bearer scan .
...
HIGH: Sensitive data stored in HTML local storage detected. [CWE-312]
https://docs.bearer.com/reference/rules/javascript_lang_session
To skip this rule, use the flag --skip-rule=javascript_lang_session

File: juice-shop/frontend/src/app/login/login.component.ts:102

 102       localStorage.setItem('email', this.user.email)


=====================================

58 checks, 38 findings

CRITICAL: 15 (CWE-22, CWE-798, CWE-89)
HIGH: 23 (CWE-312, CWE-327, CWE-548)
MEDIUM: 0
LOW: 0
WARNING: 0

exit status 1
```

Note that if only warning-level violations are found, the report does not return an exit status of 1.

```txt
24 checks, 2 warnings

CRITICAL: 0
HIGH: 0
MEDIUM: 0
LOW: 0
WARNING: 2

```

The security report is [currently available](/reference/supported-languages/) for Ruby, JavaScript, and Java projects, more languages to follow.

## Privacy Report

- Usage: `bearer scan . --report privacy`
- Default format: `csv`

The privacy report provides useful information about how your codebase uses sensitive data, with an emphasis on [data subjects](https://ico.org.uk/for-organisations/sme-web-hub/key-data-protection-terms-you-need-to-know/#datasubject) and third parties services.

The data subjects portion displays information about each detected subject and any data types associated with it. It also provides statistics on total detections and any counts of rule findings associated with the data type. In the example below, the report detects 14 instances of user telephone numbers with no rule findings.

_Note: These examples use JSON for readability, but the default format for the privacy report is CSV._

```json
"Subjects": [
  {
    "subject_name": "User",
    "name": "Telephone Number",
    "detection_count": 14,
    "critical_risk_finding_count": 0,
    "high_risk_finding_count": 0,
    "medium_risk_finding_count": 0,
    "low_risk_finding_count": 0,
    "rules_passed_count": 11
  }
]
```


The third parties portion displays data subjects and data types that are sent to or processed by known third-party services. In the example below, Bearer CLI detects a user email address sent to Sentry via the Sentry SDK and notes that a critical-risk-level rule has triggered associated with this data point.

```json
"ThirdParty": [
  {
    "third_party": "Sentry",
    "subject_name": "User",
    "data_types": [
      "Email Address"
    ],
    "critical_risk_finding_count": 1,
    "high_risk_finding_count": 0,
    "medium_risk_finding_count": 0,
    "low_risk_finding_count": 0,
    "rules_passed_count": 0
  }
]
```

The detection of third-party services is performed through an internal database knowns as Recipes. You can easily [contribute to new recipes](/contributing/recipes/).

### Customizing data subjects

By default, Bearer CLI maps all subjects to “User”, but you can override this by supplying Bearer CLI with custom mappings. This is done by passing the path to a JSON file with the `--data-subject-mapping` flag when you run the privacy report. For example:

```bash
bearer scan . --report=privacy --data-subject-mapping=/path/to/mappings.json
```

The custom map file should follow the format used by [subject_mapping.json]({{meta.sourcePath}}/blob/main/internal/classification/db/subject_mapping.json). Replace a key’s value with the higher-level subject you’d like to associate it with. Some examples might include Customer, Employee, Client, Patient, etc. Bearer CLI will use your replacement file instead of the default, so make sure to include any and all subjects you want reported.

## Data Flow Report

- Usage: `bearer scan . --report dataflow`
- Default format: `json`

The data flow report breaks down the data types and associated components detected in your code. It highlights areas in your code that process personal and sensitive data and where this data may be exposed to third parties and databases.

You can use this to gain more detailed insights beyond what the Privacy report offers, and build additional documentation like data catalogs. In the following example, we can see all of the places the [Bear Publishing](https://github.com/Bearer/bear-publishing) app processes `Email Address`:

```json
{
  "data_types": [
    {
      "name": "Email Address",
      "detectors": [
        {
          "name": "ruby",
          "locations": [
            {
              "filename": "app/controllers/application_controller.rb",
              "line_number": 35,
              "field_name": "email",
              "object_name": "current_user",
              "subject_name": "User"
            },
            {
              "filename": "app/controllers/application_controller.rb",
              "line_number": 37,
              "field_name": "email",
              "object_name": "current_user",
              "subject_name": "User"
            },
            ...
          ]
        },
        {
          "name": "schema_rb",
          "locations": [
            {
              "filename": "db/schema.rb",
              "line_number": 91,
              "stored": true,
              "parent": {
                ...
              },
              "field_name": "email",
              "object_name": "users",
              "subject_name": "User"
            }
          ]
        }
      ]
    },
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

## Next steps

For additional options on generating reports, selecting format types, and writing the output to a file, see the [command reference](/reference/commands/) documentation.
