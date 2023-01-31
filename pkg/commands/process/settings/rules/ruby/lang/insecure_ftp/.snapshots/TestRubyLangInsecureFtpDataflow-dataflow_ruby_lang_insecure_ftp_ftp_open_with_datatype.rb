data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_open_with_datatype.rb
              line_number: 6
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_insecure_ftp
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_open_with_datatype.rb
              line_number: 6
              parent:
                line_number: 3
                content: |-
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
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: ruby_lang_insecure_ftp
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_open_with_datatype.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
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
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
        - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_open_with_datatype.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
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
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
        - filename: pkg/commands/process/settings/rules/ruby/lang/insecure_ftp/testdata/ftp_open_with_datatype.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
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
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
components: []


--

