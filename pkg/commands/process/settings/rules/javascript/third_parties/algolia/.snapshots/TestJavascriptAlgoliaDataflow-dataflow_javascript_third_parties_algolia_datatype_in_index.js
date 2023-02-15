data_types:
    - name: Bank Account
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_index.js
              line_number: 4
              field_name: bank_account_number
              object_name: company
risks:
    - detector_id: javascript_third_parties_algolia
      data_types:
        - name: Bank Account
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_index.js
              line_number: 4
              parent:
                line_number: 4
                content: myAlgolia.initIndex(company.bank_account_number)
              field_name: bank_account_number
              object_name: company
    - detector_id: javascript_third_parties_algolia_client
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_index.js
          line_number: 2
          parent:
            line_number: 2
            content: algoliasearch("123", "123")
          content: |
            $<MODULE>($<_>, $<_>)
    - detector_id: javascript_third_parties_algolia_index
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_index.js
          line_number: 4
          parent:
            line_number: 4
            content: myAlgolia.initIndex(company.bank_account_number)
          content: |
            $<CLIENT>.initIndex()
components: []


--

