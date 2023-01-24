data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
              line_number: 3
              field_name: email
              object_name: user
    - name: Physical Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
              line_number: 4
              field_name: address
              object_name: user
risks:
    - detector_id: ruby_lang_logger
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    logger.info(
                      "user info:",
                      user.email,
                      user.address
                    )
              field_name: email
              object_name: user
        - name: Physical Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
              line_number: 4
              parent:
                line_number: 1
                content: |-
                    logger.info(
                      "user info:",
                      user.email,
                      user.address
                    )
              field_name: address
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    logger.info(
                      "user info:",
                      user.email,
                      user.address
                    )
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/logger/testdata/datatype_leak.rb
              line_number: 4
              parent:
                line_number: 1
                content: |-
                    logger.info(
                      "user info:",
                      user.email,
                      user.address
                    )
              object_name: user
components: []


--

