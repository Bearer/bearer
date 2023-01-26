data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/verify_none.rb
              line_number: 3
              field_name: email
              object_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/verify_none.rb
              line_number: 3
              field_name: name
              object_name: User
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/verify_none.rb
              line_number: 3
              field_name: password
              object_name: User
risks:
    - detector_id: ruby_rails_insecure_smtp
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_smtp/testdata/verify_none.rb
          line_number: 8
          parent:
            line_number: 8
            content: 'openssl_verify_mode: "none"'
          content: |
            Rails.application.configure do
              config.action_mailer.smtp_settings = {
                $<!>openssl_verify_mode: "none"
              }
            end
components: []


--

