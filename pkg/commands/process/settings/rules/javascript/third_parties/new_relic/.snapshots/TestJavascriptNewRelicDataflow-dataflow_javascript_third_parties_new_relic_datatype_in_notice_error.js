data_types:
    - name: IP address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_notice_error.js
              line_number: 7
              field_name: ip_address
              object_name: customer
              subject_name: User
risks:
    - detector_id: javascript_third_parties_new_relic
      data_types:
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/third_parties/new_relic/testdata/datatype_in_notice_error.js
              line_number: 7
              parent:
                line_number: 7
                content: newrelic.noticeError(err, customer.ip_address)
              field_name: ip_address
              object_name: customer
              subject_name: User
components: []


--

