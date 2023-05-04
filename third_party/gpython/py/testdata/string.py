# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from libtest import assertRaisesText, assertRaises

doc="format"

assert "{:#>6}".format("hello") == "#hello"
assert "{:*^6}".format(34) == "**34**"
# 填充用法
assert "{:6s}".format("hello") == " hello"
assert "{:<6s}".format("hello") == "hello "
assert "{:>6s}".format("hello") == " hello"

assert "{:*^6s}".format("hello") == "hello*"
assert "{:*^7s}".format("hello") == "*hello*"
assert "{:*^6d}".format(34) == "**34**"
assert "{:*^8s}".format("hello") == "*hello**"

assert "{.test}{[test]}".format({'test': 123}, {'test': 456}) == "123456"
assert "{[0]}{0[1]}".format([123, 456]) == "123456"

assert "{1.test}{0.test}".format({'test': 123}, {'test': 456}) == "456123"
assert "{obj[test]}".format(obj={'test': 123}) == "123"
assert "{obj.test}".format(obj={'test': 123}) == "123"
assert "{obj[val].test}".format(obj={'val': {'test': 123}}) == "123"
assert "{obj.val.test}".format(obj={'val': {'test': 123}}) == "123"

class obj:
    val = 123
assert "{obj.val}".format(obj=obj()) == "123"

# 基本用法
assert "{} {}".format("hello", "world") == "hello world"
assert "{1} {0}".format("world", "hello") == "hello world"

# 带格式化的用法
assert "{:6s}".format('aa') == "    aa"
assert "{:.2f}".format(3.14159) == "3.14"
assert "{:+d}".format(42) == "+42"
assert "{:06d}".format(42) == "000042"
assert "{:6d}".format(42) == "    42"
assert "{:d}".format(-42) == "-42"
assert "{:#o}".format(42) == "0o52"
assert "{:#x}".format(42) == "0x2a"
assert "{:#X}".format(42) == "0X2A"

# 字典访问用法
person = {"name": "Alice", "age": 30}
assert "Name: {p[name]}; Age: {p[age]}".format(p=person) == "Name: Alice; Age: 30"
assert "Name: {p[name]}; Age: {p[age]:d}".format(p=person) == "Name: Alice; Age: 30"
assert "Name: {p[name]}; Age: {p[age]:5d}".format(p=person) == "Name: Alice; Age:    30"

# 索引和字典访问混合用法
data = {"person": person, "greeting": "Hello"}
assert "{0[greeting]}, {0[person][name]}! You are {0[person][age]} years old.".format(data) == "Hello, Alice! You are 30 years old."
assert "{data[person][name]}'s age is {data[person][age]}.".format(data=data) == "Alice's age is 30."

# 指定参数名
assert "{greeting} {name}!".format(greeting="hello", name="world") == "hello world!"

# 位置参数和关键字参数混用
assert "{0} {greeting} {name} {1}".format("say", "to", greeting="hello", name="world") == "say hello world to"

# 替换大括号 FIXME
# assert "{{Hello}} {name}".format(name="World") == "{Hello} World"

try:
    "{} {} {}".format("hello", "world")
except IndexError:
    pass
else:
    assert False, "IndexError not raised"

try:
    "{0} {2}".format("hello", "world")
except IndexError:
    pass
else:
    assert False, "IndexError not raised"

try:
    "{0}".format()
except IndexError:
    pass
else:
    assert False, "IndexError not raised"

try:
    "{0.1}".format(0)
except AttributeError:
    pass
else:
    assert False, "AttributeError not raised"

try:
    "{0.key}".format(0)
except AttributeError:
    pass
else:
    assert False, "AttributeError not raised"

try:
    "{0[1]}".format(0)
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

try:
    "{0[key]}".format(0)
except TypeError:
    pass
else:
    assert False, "TypeError not raised"

# 浮点数用法
assert "{:.2f}".format(3.14159) == "3.14"
assert "{:.4f}".format(0.12345) == "0.1235"
assert "{:10.3f}".format(12.34567) == "    12.346"
assert "{:+.2f}".format(3.14) == "+3.14"
assert "{:.2e}".format(12345.0) == "1.23e+04"
assert "{:.3E}".format(12345.0) == "1.234E+04"
assert "{:10.2f}".format(123.456) == "    123.46"
assert "{:<10.2f}".format(123.456) == "123.46    "
assert "{:^10.2f}".format(123.456) == "  123.46  "
assert "{:*^10.2f}".format(123.456) == "**123.46**"

# 异常场景
try:
    "{:.2f}".format("hello")
except ValueError:
    pass
else:
    assert False, "ValueError not raised"

try:
    "{:.2e}".format("world")
except ValueError:
    pass
else:
    assert False, "ValueError not raised"

# 索引和字典访问混合
data = ["apple", "banana", {"name": "cherry"}]
assert "I like {0} and {2[name]}.".format(*data) == "I like apple and cherry."

# 混合使用
person = {"name": "Alice", "age": 25}
person_new = {"name": "Bob", "age": 30}
assert "{0.name}'s age is {0.age} and {1[name]}'s age is {1[age]}.".format(person_new, person) == "Bob's age is 30 and Alice's age is 25."

# 位置、索引和格式化混合使用
data = [
    {"name": "Alice", "age": 25},
    {"name": "Bob", "age": 30},
    {"name": "Charlie", "age": 35}
]
assert "The names are: {0[0][name]}, {0[1][name]}, {0[2][name]}.".format(data) == "The names are: Alice, Bob, Charlie."
assert "{0[1][name]:<10} {1[age]:^5}".format(data, data[1]) == "Bob         30  "
assert "{0[2][name]:>10} {1[2][age]:<5}".format(data, data) == "   Charlie 35   "
assert "{0[0][name]:.^10} {1[1][age]:+d}".format(data, data) == "..Alice... +30"

assert "{0[0][name]}".format(data) == "Alice"
assert "{0[0].name}".format(data) == "Alice"

doc="finished"
