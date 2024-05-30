class User:
    def __init__(self, name, email=""):
        self.name = name
        self.email = email

    def lowercase_name(self):
        logging.error(self.name)
        print(self.name.lower())