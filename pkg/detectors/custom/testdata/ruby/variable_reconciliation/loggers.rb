
def saveUser(user)
    # we add this properties as they are in parent scope
    user.email.username = "anon123"
    user.sex = "male"
    user.race.secondary = "hispanic"
    admin.email.username = "admin123"
    admin.isSuperUser = true

    def logUser
        # not supported as we dont support hashes variable reconciliation
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
            first_name: "John",
            last_name: "Doe",
            address: address,
            phone_number: "555-1234-123",
        }

        user.race.primary = "caucasian"
        admin.email.domain = "domain.com"


        # custom detector matches
        logger.info("User email:", user)
        logger.info("Admin domain:", admin.email)
        logger.info("student email:", student.email)
    end

    # we ignore adjecent scopes
    def parse_users(user)
        # we ignore this occurence of user as it is not in the same scope nor parent scope
        user.health_records
    end

    return user
end





