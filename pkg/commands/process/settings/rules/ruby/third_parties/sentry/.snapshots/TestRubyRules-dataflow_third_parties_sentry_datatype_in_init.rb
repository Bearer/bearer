data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_init.rb
              line_number: 3
              field_name: email
              object_name: current_user
risks:
    - detector_id: ruby_third_parties_sentry
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_init.rb
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    Sentry.init do |config|
                      config.before_breadcrumb = lambda do |breadcrumb, hint|
                        breadcrumb.message = "Authenticated user #{current_user.email}"
                        breadcrumb
                      end
                    end
              field_name: email
              object_name: current_user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/sentry/testdata/datatype_in_init.rb
              line_number: 3
              parent:
                line_number: 1
                content: |-
                    Sentry.init do |config|
                      config.before_breadcrumb = lambda do |breadcrumb, hint|
                        breadcrumb.message = "Authenticated user #{current_user.email}"
                        breadcrumb
                      end
                    end
              object_name: current_user
components: []


--

