
def saveUser(user)
    # we add this properties as they are in parent scope
    user.address.zip
    user.first_name

    def logUser
        address = {
            zip_code: "75000",
            address: "1 rue du Marché",
                city: "Paris"
            }
        
        Employee = Struct.new(email:, first_name:, last_name:, keyword_init: true)
        
        employee = Employee.new(email: "user@example.com", first_name: "John", last_name: "Doe")

        # we add this properties as they are in the same scope
        user = {
            email: "user@example.com",
            first_name: "John",
                last_name: "Doe",
                address:,
            phone_number: "555-1234-123",
        }
        
        # custom detector matches
        logger.info("User info:", user)
        logger.info("Employee info:", employee)
        logger.info("User email:", user.email.domain)
        Rails.logger.info("User info:", user)
    end

    # we ignore adjecent scopes
    def parse_users(user)
        # we ignore this occurence of user as it is not in the same scope nor parent scope
        user.last_name
    end

    return user
end





