data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
              line_number: 12
              field_name: email
              object_name: user
              subject_name: User
    - name: IP address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
              line_number: 7
              field_name: ip_address
              object_name: user
              subject_name: User
    - name: Unique Identifier
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
              line_number: 7
              field_name: user_id
              object_name: userObj
              subject_name: User
risks:
    - detector_id: javascript_third_parties_algolia
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
              line_number: 12
              parent:
                line_number: 12
                content: 'index.saveObjects([{ email: user.email }])'
              field_name: email
              object_name: user
              subject_name: User
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
              line_number: 7
              parent:
                line_number: 8
                content: |-
                    index
                      .saveObject(userObj, { autoGenerateObjectIDIfNotExist: true })
              field_name: ip_address
              object_name: user
              subject_name: User
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
              line_number: 7
              parent:
                line_number: 8
                content: |-
                    index
                      .saveObject(userObj, { autoGenerateObjectIDIfNotExist: true })
              field_name: user_id
              object_name: userObj
              subject_name: User
    - detector_id: javascript_third_parties_algolia_client
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
          line_number: 2
          parent:
            line_number: 2
            content: algoliasearch("123", "123")
          content: |
            $<MODULE>($<_>, $<_>)
    - detector_id: javascript_third_parties_algolia_index
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
          line_number: 4
          parent:
            line_number: 4
            content: myAlgolia.initIndex("test_index")
          content: |
            $<CLIENT>.initIndex()
    - detector_id: javascript_third_parties_algolia_module
      locations:
        - filename: pkg/commands/process/settings/rules/javascript/third_parties/algolia/testdata/datatype_in_save_object.js
          line_number: 1
          parent:
            line_number: 1
            content: require("algoliasearch")
          content: |
            require("algoliasearch")
components: []


--

