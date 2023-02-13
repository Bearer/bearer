data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/fail_with_different_version.rb
              line_number: 5
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_dataflow
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/fail_with_different_version.rb
              line_number: 5
              parent:
                line_number: 5
                content: 'templates_client.create_job_from_template(project_id: "123", job_name: "my-job", parameters: { current_user: user.email })'
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: google_dataflow_client_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/fail_with_different_version.rb
          line_number: 3
          parent:
            line_number: 3
            content: Google::Cloud::Dataflow::V2::TemplatesService::Client.new
          content: |
            Google::Cloud::Dataflow::$<_>::$<_>::Client.new
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/fail_with_different_version.rb
          line_number: 4
          parent:
            line_number: 4
            content: Google::Cloud::Dataflow.templates_service
          content: |
            Google::Cloud::Dataflow.$<METHOD>
    - detector_id: google_dataflow_templates_service_client_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/fail_with_different_version.rb
          line_number: 3
          parent:
            line_number: 3
            content: Google::Cloud::Dataflow::V2::TemplatesService::Client.new
          content: |
            Google::Cloud::Dataflow::$<_>::TemplatesService::Client.new
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/fail_with_different_version.rb
          line_number: 4
          parent:
            line_number: 4
            content: Google::Cloud::Dataflow.templates_service
          content: |
            Google::Cloud::Dataflow.templates_service
components: []


--

