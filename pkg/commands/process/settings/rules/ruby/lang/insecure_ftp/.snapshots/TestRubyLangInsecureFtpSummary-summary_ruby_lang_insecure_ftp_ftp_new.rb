critical:
    - rule_dsrid: DSR-2
      rule_display_id: ruby_lang_insecure_ftp
      rule_description: Only communicate using SFTP connections.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_insecure_ftp
      line_number: 8
      filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_new.rb
      category_groups:
        - PII
      parent_line_number: 8
      parent_content: Net::FTP.new("ftp.ruby-lang.org")


--

