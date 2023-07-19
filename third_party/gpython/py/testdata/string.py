# Copyright 2018 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import libtest as self

doc="join"

class MyWrapper:
    def __init__(self, sval): self.sval = sval
    def __str__(self): return self.sval


self.assertEqual(' '.join(['a', 'b', 'c', 'd']), 'a b c d')
self.assertEqual(''.join(('a', 'b', 'c', 'd')), 'abcd')
self.assertEqual(' '.join('wxyz'), 'w x y z')
self.assertEqual(' '.join(['a', 'b', 'c', 'd']), 'a b c d')

self.assertEqual(' '.join(['a', 'b', 'c', 'd']), 'a b c d')
self.assertEqual(''.join(('a', 'b', 'c', 'd')), 'abcd')
self.assertEqual(' '.join('wxyz'), 'w x y z')

self.assertEqual(' '.join(['1', '2', MyWrapper('foo')]), '1 2 foo')
self.assertEqual(' '.join(['1', '2', '3', bytes()]), "1 2 3 b''")
self.assertEqual(' '.join([1, 2, 3]), '1 2 3')
self.assertEqual(' '.join(['1', '2', 3]), '1 2 3')

# FIXME: self.assertRaises(TypeError, ' '.join, ['1', '2', MyWrapper('foo')])
# FIXME: self.assertRaises(TypeError, ' '.join, ['1', '2', '3', bytes()])
# FIXME: self.assertRaises(TypeError, ' '.join, [1, 2, 3])
# FIXME: self.assertRaises(TypeError, ' '.join, ['1', '2', 3])

doc="format"

self.assertEqual("{:#>6}".format("hello"), "#hello")
self.assertEqual("{:*^6}".format(34), "**34**")
# 填充用法
self.assertEqual("{:6s}".format("hello"), " hello")
self.assertEqual("{:<6s}".format("hello"), "hello ")
self.assertEqual("{:>6s}".format("hello"), " hello")

self.assertEqual("{:*^6s}".format("hello"), "hello*")
self.assertEqual("{:*^7s}".format("hello"), "*hello*")
self.assertEqual("{:*^6d}".format(34), "**34**")
self.assertEqual("{:*^8s}".format("hello"), "*hello**")

self.assertEqual("{.test}{[test]}".format({'test': 123}, {'test': 456}), "123456")
self.assertEqual("{[0]}{0[1]}".format([123, 456]), "123456")

self.assertEqual("{1.test}{0.test}".format({'test': 123}, {'test': 456}), "456123")
self.assertEqual("{obj[test]}".format(obj={'test': 123}), "123")
self.assertEqual("{obj.test}".format(obj={'test': 123}), "123")
self.assertEqual("{obj[val].test}".format(obj={'val': {'test': 123}}), "123")
self.assertEqual("{obj.val.test}".format(obj={'val': {'test': 123}}), "123")

class obj:
    val = 123
self.assertEqual("{obj.val}".format(obj=obj()), "123")

# 基本用法
self.assertEqual("{} {}".format("hello", "world"), "hello world")
self.assertEqual("{1} {0}".format("world", "hello"), "hello world")

# 带格式化的用法
self.assertEqual("{:6s}".format('aa'), "    aa")
self.assertEqual("{:.2f}".format(3.14159), "3.14")
self.assertEqual("{:+d}".format(42), "+42")
self.assertEqual("{:06d}".format(42), "000042")
self.assertEqual("{:6d}".format(42), "    42")
self.assertEqual("{:d}".format(-42), "-42")
self.assertEqual("{:#o}".format(42), "0o52")
self.assertEqual("{:#x}".format(42), "0x2a")
self.assertEqual("{:#X}".format(42), "0X2A")

# 字典访问用法
person = {"name": "Alice", "age": 30}
self.assertEqual("Name: {p[name]}; Age: {p[age]}".format(p=person), "Name: Alice; Age: 30")
self.assertEqual("Name: {p[name]}; Age: {p[age]:d}".format(p=person), "Name: Alice; Age: 30")
self.assertEqual("Name: {p[name]}; Age: {p[age]:5d}".format(p=person), "Name: Alice; Age:    30")

# 索引和字典访问混合用法
data = {"person": person, "greeting": "Hello"}
self.assertEqual("{0[greeting]}, {0[person][name]}! You are {0[person][age]} years old.".format(data), "Hello, Alice! You are 30 years old.")
self.assertEqual("{data[person][name]}'s age is {data[person][age]}.".format(data=data), "Alice's age is 30.")

# 指定参数名
self.assertEqual("{greeting} {name}!".format(greeting="hello", name="world"), "hello world!")

# 位置参数和关键字参数混用
self.assertEqual("{0} {greeting} {name} {1}".format("say", "to", greeting="hello", name="world"), "say hello world to")

# 替换大括号 FIXME
# assert "{{Hello}} {name}".format(name="World") == "{Hello} World"

self.assertRaises(IndexError, "{} {} {}".format, "hello", "world")
self.assertRaises(IndexError, "{0} {2}".format, "hello", "world")
self.assertRaises(IndexError, "{0}".format)
self.assertRaises(AttributeError, "{0.1}".format, 0)
self.assertRaises(AttributeError, "{0.key}".format, 0)
self.assertRaises(TypeError, "{0[1]}".format, 0)
self.assertRaises(TypeError, "{0[key]}".format, 0)


# 浮点数用法
self.assertEqual("{:.2f}".format(3.14159), "3.14")
self.assertEqual("{:.4f}".format(0.12345), "0.1235")
self.assertEqual("{:10.3f}".format(12.34567), "    12.346")
self.assertEqual("{:+.2f}".format(3.14), "+3.14")
self.assertEqual("{:.2e}".format(12345.0), "1.23e+04")
self.assertEqual("{:.3E}".format(12345.0), "1.234E+04")
self.assertEqual("{:10.2f}".format(123.456), "    123.46")
self.assertEqual("{:<10.2f}".format(123.456), "123.46    ")
self.assertEqual("{:^10.2f}".format(123.456), "  123.46  ")
self.assertEqual("{:*^10.2f}".format(123.456), "**123.46**")

# 异常场景
self.assertRaises(ValueError, "{:.2f}".format, "hello")
self.assertRaises(ValueError, "{:.2e}".format, "world")

# 索引和字典访问混合
data = ["apple", "banana", {"name": "cherry"}]
self.assertEqual("I like {0} and {2[name]}.".format(*data), "I like apple and cherry.")

# 混合使用
person = {"name": "Alice", "age": 25}
person_new = {"name": "Bob", "age": 30}
self.assertEqual("{0.name}'s age is {0.age} and {1[name]}'s age is {1[age]}.".format(person_new, person), "Bob's age is 30 and Alice's age is 25.")

# 位置、索引和格式化混合使用
data = [
    {"name": "Alice", "age": 25},
    {"name": "Bob", "age": 30},
    {"name": "Charlie", "age": 35}
]
self.assertEqual("The names are: {0[0][name]}, {0[1][name]}, {0[2][name]}.".format(data), "The names are: Alice, Bob, Charlie.")
self.assertEqual("{0[1][name]:<10} {1[age]:^5}".format(data, data[1]), "Bob         30  ")
self.assertEqual("{0[2][name]:>10} {1[2][age]:<5}".format(data, data), "   Charlie 35   ")
self.assertEqual("{0[0][name]:.^10} {1[1][age]:+d}".format(data, data), "..Alice... +30")

self.assertEqual("{0[0][name]}".format(data), "Alice")
self.assertEqual("{0[0].name}".format(data), "Alice")

doc="finished"
