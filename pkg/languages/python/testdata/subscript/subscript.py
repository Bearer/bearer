foo[input()]

globals()[input()]

x = input()

foo[x]

globals()[x]
globals()[x]("something")

# don't match the below

foo[y]
foo["world"]
globals()[y]
globals()["world"]