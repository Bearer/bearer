critical:
    - rule_dsrid: DSR-4
      rule_display_id: javascript_lang_file_generation
      rule_description: Do not write sensitive data to static files.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_lang_file_generation
      line_number: 8
      filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
      category_groups:
        - PII
      parent_line_number: 18
      parent_content: |-
        fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
          if (err) console.log(err)
          else console.log("Data saved")
        })
    - rule_dsrid: DSR-4
      rule_display_id: javascript_lang_file_generation
      rule_description: Do not write sensitive data to static files.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_lang_file_generation
      line_number: 11
      filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
      category_groups:
        - PII
      parent_line_number: 18
      parent_content: |-
        fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
          if (err) console.log(err)
          else console.log("Data saved")
        })
    - rule_dsrid: DSR-4
      rule_display_id: javascript_lang_file_generation
      rule_description: Do not write sensitive data to static files.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_lang_file_generation
      line_number: 12
      filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
      category_groups:
        - PII
      parent_line_number: 18
      parent_content: |-
        fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
          if (err) console.log(err)
          else console.log("Data saved")
        })


--

