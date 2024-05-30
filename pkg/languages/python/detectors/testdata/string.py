class Greet:
    Greeting = "Hello World"

    def main(args):
        s = Greet.Greeting + "!"
        s += "!!"

        s2 = "hey "
        s2 += args[0]
        s2 += " there"

        s3 = f"foo '{s2}' bar"
