risks:
    - detector_id: ruby_rails_insecure_communication
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_communication/testdata/no_datatypes.rb
          line_number: 2
          parent:
            line_number: 2
            content: config.force_ssl = false
          content: |
            Rails.application.configure do
              $<!>config.force_ssl = false
            end
components: []


--

