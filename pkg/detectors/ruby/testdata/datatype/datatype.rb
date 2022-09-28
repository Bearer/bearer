class Shape
	attr_accessor :name
	def initialize(name, color)
		@name = name
		@name[:surname] = color
	end
end

@person[:city][:number]
@person.city.number

person[:city][:number]
person.city.number

person[:city].test[:number].test()
@person[:city].test[:number].test()


person[:city].test()[:number].test()
@person[:city].test()[:number].test()


person.test().city.address().test[:city].address()
@person[:city].test()[:number].test()

File.join

Foo = Struct.new(foo: "foo", bar: "bar")
Foo.new(foo: "foo", bar: "bar")


Foo = Class.new(Base) do
	attr_accessor :name
	def initialize()
		@name[:surname] = color
	end
end