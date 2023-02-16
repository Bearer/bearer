data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_set_custom_attribute.js
              line_number: 3
              field_name: email
              object_name: customer
              subject_name: User
risks:
    - detector_id: javascript_third_parties_new_relic
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_set_custom_attribute.js
              line_number: 3
              parent:
                line_number: 3
                content: newrelic.setCustomAttribute("user-id", customer.email)
              field_name: email
              object_name: customer
              subject_name: User
components: []


--

