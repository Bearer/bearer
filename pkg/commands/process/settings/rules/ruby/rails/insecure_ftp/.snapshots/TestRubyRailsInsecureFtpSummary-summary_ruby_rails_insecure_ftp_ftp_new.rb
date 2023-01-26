critical:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_rails_insecure_ftp
      policy_description: Only communicate using SFTP connections.
      line_number: 8
      filename: pkg/commands/process/settings/rules/ruby/rails/insecure_ftp/testdata/ftp_new.rb
      category_groups:
        - PII
      parent_line_number: 8
      parent_content: Net::FTP.new("ftp.ruby-lang.org")


--

