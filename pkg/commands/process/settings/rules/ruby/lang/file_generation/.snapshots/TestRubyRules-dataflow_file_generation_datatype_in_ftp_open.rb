data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 1
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: email
              object_name: user
    - name: Emails
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 3
              field_name: emails
              object_name: user
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: first_name
              object_name: user
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: last_name
              object_name: user
risks:
    - detector_id: ruby_lang_file_generation
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'File.open("users.log", "w") { |f| f.write "#{Time.now} - User #{user.email} logged in\n" }'
              field_name: email
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user.emails, "users.csv", "w") do |f|
                    	users.each do |user|
                    		f.write "#{user.email},#{user.first_name},#{user.last_name}"
                    	end
                    end
              field_name: email
              object_name: user
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user.emails, "users.csv", "w") do |f|
                    	users.each do |user|
                    		f.write "#{user.email},#{user.first_name},#{user.last_name}"
                    	end
                    end
              field_name: first_name
              object_name: user
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user.emails, "users.csv", "w") do |f|
                    	users.each do |user|
                    		f.write "#{user.email},#{user.first_name},#{user.last_name}"
                    	end
                    end
              field_name: last_name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'File.open("users.log", "w") { |f| f.write "#{Time.now} - User #{user.email} logged in\n" }'
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user.emails, "users.csv", "w") do |f|
                    	users.each do |user|
                    		f.write "#{user.email},#{user.first_name},#{user.last_name}"
                    	end
                    end
              object_name: user
components: []


--

