critical:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_insecure_ftp
      policy_description: Only communicate using SFTP connections.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_open_with_datatype.rb
      category_groups:
        - PII
      parent_line_number: 3
      parent_content: |-
        Net::FTP.open("ftp.site.com") do |ftp|
          file = Tempfile.new("user_data")
          begin
            file << user.email
            file.close

            ftp.puttextfile(file.path, "/users/123.json")
          ensure
            file.close!
          end
        end


--

