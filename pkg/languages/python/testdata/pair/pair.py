user_input = input("Enter username: ")

# collection is some mongodo collection
collection.find_one({"username": user_input})