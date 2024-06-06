from baz import foo
foo.someMethod()

from baz import foo as asdf
asdf.someMethod()

from baz import y as z, a as b, foo as j
j.someMethod()

from baz import y, a, foo
foo.someMethod()

import bar
bar.someMethod()

import xyz, bar
bar.someMethod()

import bar as qwerty
qwerty.someMethod()

import yy as zz, bar as bb
bb.someMethod()

import foo.bat
foo.bat.dottedMethod()

import FooClass
z = FooClass
z.qwerty()

from baz import FooClass as Something
x = Something()
x.qwerty()

import FooClass as SomethingElse
y = SomethingElse()
y.qwerty()

foo()