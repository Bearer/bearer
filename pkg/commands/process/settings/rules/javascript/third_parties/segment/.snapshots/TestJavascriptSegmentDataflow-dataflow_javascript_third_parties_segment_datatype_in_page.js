data_types:
    - name: Bank Account
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_page.js
              line_number: 10
              field_name: bank_account_number
              object_name: customer
              subject_name: User
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_page.js
              line_number: 7
              field_name: userId
              object_name: page
risks:
    - detector_id: javascript_third_parties_segment
      data_types:
        - name: Bank Account
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_page.js
              line_number: 10
              parent:
                line_number: 6
                content: |-
                    analytics.page({
                      userId: customer.id,
                      category: "Shopping Cart",
                      properties: {
                        path: "/cart/"+customer.bank_account_number
                      },
                    })
              field_name: bank_account_number
              object_name: customer
              subject_name: User
    - detector_id: javascript_third_parties_segment_instance
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/segment/testdata/datatype_in_page.js
          line_number: 3
          parent:
            line_number: 3
            content: 'new Analytics({ writeKey: ''my-write-key'' })'
          content: |
            new Analytics()
components: []


--

