critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_bugsnag
      rule_description: Do not send sensitive data to Bugsnag.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_bugsnag
      line_number: 3
      filename: pkg/commands/process/settings/rules/javascript/third_parties/bugsnag/testdata/datatype_in_start.js
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Bugsnag.start({
          onError: function (e) {
            e.setUser(user.id, user.email, user.name)
            e.addMetadata('user location', {
              country: user.home_country,
            })
          },
          onSession: function (session) {
            session.setUser(user.email)
          }
        })
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_bugsnag
      rule_description: Do not send sensitive data to Bugsnag.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_bugsnag
      line_number: 5
      filename: pkg/commands/process/settings/rules/javascript/third_parties/bugsnag/testdata/datatype_in_start.js
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Bugsnag.start({
          onError: function (e) {
            e.setUser(user.id, user.email, user.name)
            e.addMetadata('user location', {
              country: user.home_country,
            })
          },
          onSession: function (session) {
            session.setUser(user.email)
          }
        })
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_bugsnag
      rule_description: Do not send sensitive data to Bugsnag.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_bugsnag
      line_number: 9
      filename: pkg/commands/process/settings/rules/javascript/third_parties/bugsnag/testdata/datatype_in_start.js
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        Bugsnag.start({
          onError: function (e) {
            e.setUser(user.id, user.email, user.name)
            e.addMetadata('user location', {
              country: user.home_country,
            })
          },
          onSession: function (session) {
            session.setUser(user.email)
          }
        })


--

