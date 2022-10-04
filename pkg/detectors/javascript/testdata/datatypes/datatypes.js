let user;
user.Save();
user.name;
user.name = "test";
user.name.surname;

user.name.test.surname.Save();

// user.test is one chain while name.surname is another
user.test["test"].name.surname;

// 3 chains ->
//  user->test
//  test->something
//  todo->list
user.test((test) => test.something).todo.list;

let test = {
	address: {
		city: "new york",
		street: {
			letters: "11-th avenue",
			numer: 12,
		},
	},
};

// vue.config.js
const path = require("path");
const BundleAnalyzerPlugin =
	require("webpack-bundle-analyzer").BundleAnalyzerPlugin;

function resolve(dir) {
	return path.join(__dirname, dir);
}

this.save();

/// scoping datatypes

// scope program
myobject.Method();
myobject.Method1();

const a = () => {
	// scope arrow_function
	myobject1.Method();
	myobject1.Method1();
};

function doSomething() {
	// scope function
	myobject2.Method();
	myobject2.Method1();
}

var a = {
	doSomething() {
		// scope method_definition
		myobject3.Method();
		myobject3.Method1();
	},
};

a.test.something = {
	name: "user",
	5: {
		test: "test",
	},
};

a = {
	test: "test",
};

test[test.name] = {
	el: "test",
};

function test() {
	return { name: "name" };
}

const test = () => ({ name: "name" });

a.test = () => ({ name: "name" });

a["test"] = () => ({ name: "name" });

module.exports = function (url) {
	return {
		path: "user",
	};
};

class User {
	constructor(name, age) {
		this.name = name;
		this.age = age;
	}

	get status() {
		return `${this.name} -> ${this.age}`;
	}
}

var user = new User("cedric", "35");

console.log(user.status);
