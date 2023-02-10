data_types:
    - name: Email Address
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 11
              field_name: email
              object_name: user
              subject_name: User
    - name: Firstname
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 8
              field_name: firstname
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 12
              field_name: firstname
              object_name: user
              subject_name: User
    - name: Lastname
      detectors:
        - name: javascript
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 8
              field_name: surname
              object_name: user
              subject_name: User
risks:
    - detector_id: javascript_lang_file_generation
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 11
              parent:
                line_number: 18
                content: |-
                    fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
                      if (err) console.log(err)
                      else console.log("Data saved")
                    })
              field_name: email
              object_name: user
              subject_name: User
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 8
              parent:
                line_number: 18
                content: |-
                    fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
                      if (err) console.log(err)
                      else console.log("Data saved")
                    })
              field_name: firstname
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 12
              parent:
                line_number: 18
                content: |-
                    fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
                      if (err) console.log(err)
                      else console.log("Data saved")
                    })
              field_name: firstname
              object_name: user
              subject_name: User
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/file_generation/testdata/file_generation.js
              line_number: 8
              parent:
                line_number: 18
                content: |-
                    fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
                      if (err) console.log(err)
                      else console.log("Data saved")
                    })
              field_name: surname
              object_name: user
              subject_name: User
components: []


--

