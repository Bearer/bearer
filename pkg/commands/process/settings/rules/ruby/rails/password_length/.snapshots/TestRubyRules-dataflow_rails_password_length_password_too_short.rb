data_types:
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/password_length/testdata/password_too_short.rb
              line_number: 3
              field_name: password
              object_name: User
risks:
    - detector_id: ruby_rails_password_length
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/rails/password_length/testdata/password_too_short.rb
          line_number: 3
          parent:
            line_number: 3
            content: 'validates :password, length: { minimum: 6 }'
          content: |
            class $<_>
              $<!>validates :password, length: { minimum: $<LENGTH> }
            end
components: []


--

