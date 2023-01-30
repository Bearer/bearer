data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/ssl_disabled.rb
              line_number: 3
              field_name: email
              object_name: User
              subject_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/ssl_disabled.rb
              line_number: 3
              field_name: name
              object_name: User
              subject_name: User
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/ssl_disabled.rb
              line_number: 3
              field_name: password
              object_name: User
              subject_name: User
risks:
    - detector_id: ruby_rails_insecure_communication
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/ssl_disabled.rb
          line_number: 7
          parent:
            line_number: 7
            content: config.force_ssl = false
          content: |
            Rails.application.configure do
              $<!>config.force_ssl = false
            end
components: []


--

