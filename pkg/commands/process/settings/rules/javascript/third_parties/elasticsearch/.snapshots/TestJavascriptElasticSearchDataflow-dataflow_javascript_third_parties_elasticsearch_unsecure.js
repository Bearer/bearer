data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/elasticsearch/testdata/unsecure.js
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
    - name: Fullname
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/elasticsearch/testdata/unsecure.js
              line_number: 1
              field_name: name
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_elasticsearch
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/elasticsearch/testdata/unsecure.js
              line_number: 1
              parent:
                line_number: 2
                content: elasticsearch.index(user)
              field_name: email
              object_name: user
              subject_name: User
        - name: Fullname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/elasticsearch/testdata/unsecure.js
              line_number: 1
              parent:
                line_number: 2
                content: elasticsearch.index(user)
              field_name: name
              object_name: user
              subject_name: User
components: []


--

