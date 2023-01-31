data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/logger/testdata/datatype_leak.rb
              line_number: 1
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_rails_logger
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/logger/testdata/datatype_leak.rb
              line_number: 1
              parent:
                line_number: 1
                content: Rails.logger.info(user.email)
              field_name: email
              object_name: user
              subject_name: User
components: []


--

