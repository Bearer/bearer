risks:
    - detector_id: javascript_jwt
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/javascript/lang/jwt/testdata/unsecure.js
              line_number: 2
              parent:
                line_number: 2
                content: 'jwt.sign({ user: { email: "jhon@gmail.com" } }, "shhhhh")'
              field_name: email
              object_name: user
              subject_name: User
components: []


--

