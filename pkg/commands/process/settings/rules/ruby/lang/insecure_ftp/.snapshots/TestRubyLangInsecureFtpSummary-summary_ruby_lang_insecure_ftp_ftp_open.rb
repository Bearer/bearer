low:
    - policy_name: ""
      policy_dsrid: DSR-2
      policy_display_id: ruby_lang_insecure_ftp
      policy_description: Only communicate using SFTP connections.
      line_number: 3
      filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_open.rb
      parent_line_number: 3
      parent_content: |-
        Net::FTP.open("example.com") do |ftp|
          ftp.login
          files = ftp.chdir('pub/lang/ruby/contrib')
          files = ftp.list('n*')
          ftp.getbinaryfile('nif.rb-0.91.gz', 'nif.gz', 1024)
        end


--

