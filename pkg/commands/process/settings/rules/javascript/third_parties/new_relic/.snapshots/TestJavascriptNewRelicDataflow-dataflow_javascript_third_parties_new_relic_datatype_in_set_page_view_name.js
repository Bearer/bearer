data_types:
    - name: Bank Account
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_set_page_view_name.js
              line_number: 3
              field_name: bank_account_number
              object_name: customer
              subject_name: User
risks:
    - detector_id: javascript_third_parties_new_relic
      data_types:
        - name: Bank Account
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_set_page_view_name.js
              line_number: 3
              parent:
                line_number: 3
                content: newrelic.setPageViewName(customer.bank_account_number, "$host")
              field_name: bank_account_number
              object_name: customer
              subject_name: User
components: []


--

