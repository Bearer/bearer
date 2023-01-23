data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 1
              field_name: email
              object_name: user_2
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: email
              object_name: user_4
    - name: Emails
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 3
              field_name: emails
              object_name: user_3
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: first_name
              object_name: user_4
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              field_name: last_name
              object_name: user_4
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
                content: 'File.open("users.log", "w") { |f| f.write "#{Time.now} - User #{user_2.email} logged in\n" }'
              field_name: email
              object_name: user_2
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user_3.emails, "users.csv", "w") do |f|
                    	users.each do |user_4|
                    		f.write "#{user_4.email},#{user_4.first_name},#{user_4.last_name}"
                    	end
                    end
              field_name: email
              object_name: user_4
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user_3.emails, "users.csv", "w") do |f|
                    	users.each do |user_4|
                    		f.write "#{user_4.email},#{user_4.first_name},#{user_4.last_name}"
                    	end
                    end
              field_name: first_name
              object_name: user_4
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user_3.emails, "users.csv", "w") do |f|
                    	users.each do |user_4|
                    		f.write "#{user_4.email},#{user_4.first_name},#{user_4.last_name}"
                    	end
                    end
              field_name: last_name
              object_name: user_4
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 1
              parent:
                line_number: 1
                content: 'File.open("users.log", "w") { |f| f.write "#{Time.now} - User #{user_2.email} logged in\n" }'
              object_name: user_2
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_ftp_open.rb
              line_number: 5
              parent:
                line_number: 3
                content: |-
                    File.open(user_3.emails, "users.csv", "w") do |f|
                    	users.each do |user_4|
                    		f.write "#{user_4.email},#{user_4.first_name},#{user_4.last_name}"
                    	end
                    end
              object_name: user_4
components: []


--

