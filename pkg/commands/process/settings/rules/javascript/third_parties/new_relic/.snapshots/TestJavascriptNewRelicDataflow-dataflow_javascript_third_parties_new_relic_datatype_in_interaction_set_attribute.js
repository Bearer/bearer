data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
              line_number: 13
              field_name: email_address
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
              line_number: 6
              field_name: first_name
              object_name: user
              subject_name: User
    - name: Interactions
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
              line_number: 7
              field_name: post_code
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_third_parties_new_relic
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
              line_number: 13
              parent:
                line_number: 13
                content: interaction.setAttribute("email", user.email_address)
              field_name: email_address
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
              line_number: 6
              parent:
                line_number: 5
                content: |-
                    newrelic.interaction()
                        .setAttribute("username", user.first_name)
              field_name: first_name
              object_name: user
              subject_name: User
        - name: Interactions
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
              line_number: 7
              parent:
                line_number: 5
                content: |-
                    newrelic.interaction()
                        .setAttribute("username", user.first_name)
                        .setAttribute("postal-code", user.post_code)
              field_name: post_code
              object_name: user
              subject_name: User
    - detector_id: javascript_new_relic_interaction
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
          line_number: 5
          parent:
            line_number: 5
            content: newrelic.interaction()
          content: |
            $<CLIENT>.interaction()
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_interaction_set_attribute.js
          line_number: 12
          parent:
            line_number: 12
            content: newrelic.interaction()
          content: |
            $<CLIENT>.interaction()
components: []


--

