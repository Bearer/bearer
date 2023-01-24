risks:
    - detector_id: ruby_rails_insecure_ftp
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_ftp/testdata/ftp_open.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
                Net::FTP.open("example.com") do |ftp|
                  ftp.login
                  files = ftp.chdir('pub/lang/ruby/contrib')
                  files = ftp.list('n*')
                  ftp.getbinaryfile('nif.rb-0.91.gz', 'nif.gz', 1024)
                end
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_ftp/testdata/ftp_open.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
                Net::FTP.open("example.com") do |ftp|
                  ftp.login
                  files = ftp.chdir('pub/lang/ruby/contrib')
                  files = ftp.list('n*')
                  ftp.getbinaryfile('nif.rb-0.91.gz', 'nif.gz', 1024)
                end
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_ftp/testdata/ftp_open.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
                Net::FTP.open("example.com") do |ftp|
                  ftp.login
                  files = ftp.chdir('pub/lang/ruby/contrib')
                  files = ftp.list('n*')
                  ftp.getbinaryfile('nif.rb-0.91.gz', 'nif.gz', 1024)
                end
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_ftp/testdata/ftp_open.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
                Net::FTP.open("example.com") do |ftp|
                  ftp.login
                  files = ftp.chdir('pub/lang/ruby/contrib')
                  files = ftp.list('n*')
                  ftp.getbinaryfile('nif.rb-0.91.gz', 'nif.gz', 1024)
                end
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
        - filename: pkg/commands/process/settings/rules/ruby/rails/insecure_ftp/testdata/ftp_open.rb
          line_number: 3
          parent:
            line_number: 3
            content: |-
                Net::FTP.open("example.com") do |ftp|
                  ftp.login
                  files = ftp.chdir('pub/lang/ruby/contrib')
                  files = ftp.list('n*')
                  ftp.getbinaryfile('nif.rb-0.91.gz', 'nif.gz', 1024)
                end
          content: |
            $<!>Net::FTP.open() do
              $<_>
            end
components: []


--

