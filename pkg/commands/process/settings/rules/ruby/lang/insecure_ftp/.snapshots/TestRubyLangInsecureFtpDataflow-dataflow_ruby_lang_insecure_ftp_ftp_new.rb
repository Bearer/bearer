data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_new.rb
              line_number: 5
              field_name: email
              object_name: User
    - name: Fullname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_new.rb
              line_number: 5
              field_name: name
              object_name: User
    - name: Passwords
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_new.rb
              line_number: 5
              field_name: password
              object_name: User
risks:
    - detector_id: ruby_lang_insecure_ftp
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_new.rb
          line_number: 8
          parent:
            line_number: 8
            content: Net::FTP.new("ftp.ruby-lang.org")
          content: |
            Net::FTP.new()
components: []


--

