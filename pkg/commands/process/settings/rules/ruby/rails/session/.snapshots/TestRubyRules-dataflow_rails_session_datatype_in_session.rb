data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/session/testdata/datatype_in_session.rb
              line_number: 1
              field_name: email
              object_name: user
    - name: Unique Identifier
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/session/testdata/datatype_in_session.rb
              line_number: 2
              field_name: user_id
              object_name: session
risks:
    - detector_id: ruby_rails_session
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/session/testdata/datatype_in_session.rb
              line_number: 1
              parent:
                line_number: 1
                content: session[:current_user] = user.email
              field_name: email
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/rails/session/testdata/datatype_in_session.rb
              line_number: 1
              parent:
                line_number: 1
                content: session[:current_user] = user.email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/rails/session/testdata/datatype_in_session.rb
              line_number: 2
              parent:
                line_number: 2
                content: session[:user_id] = current_user.user_id
              field_name: user_id
              object_name: current_user
components: []


--

