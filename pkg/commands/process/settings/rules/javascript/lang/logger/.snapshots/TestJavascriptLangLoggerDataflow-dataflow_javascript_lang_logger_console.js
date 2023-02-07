data_types:
    - name: Fullname
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/console.js
              line_number: 1
              field_name: name
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_lang_logger
      data_types:
        - name: Fullname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/logger/testdata/console.js
              line_number: 1
              parent:
                line_number: 1
                content: console.log(user.name)
              field_name: name
              object_name: user
              subject_name: User
components: []


--

