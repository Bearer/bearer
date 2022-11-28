
def saveUser(user)
    # we add this properties as they are in parent scope
    user.email.username
    user.race

    def logUser
        address = {
            zip_code: "75000",
            address: "1 rue du Marché",
                city: "Paris"
            }
        
        # this kind of structure definitions are not supported at the moment
        Employee = Struct.new(email:, first_name:, last_name:, keyword_init: true)
        employee = Employee.new(email: "user@example.com", first_name: "John", last_name: "Doe")

        # we add this properties as they are in the same scope
        user = {
            email: "user@example.com",
            first_name: "John",
            last_name: "Doe",
            address: address,
            phone_number: "555-1234-123",
        }
        
        # custom detector ignores, it ignores this ones as they don't have any properties tied to them
        logger.info("User info:", user)
        logger.info("Employee info:", employee)
        Rails.logger.info("User info:", user)

        # custom detector matches
        logger.info("User email:", user.email.domain)
    end

    # we ignore adjecent scopes
    def parse_users(user)
        # we ignore this occurence of user as it is not in the same scope nor parent scope
        user.health_records
    end

    return user
end





