import re

print("map test")
ls = ["  test1", "aaa", "test2 ", " test3 ", " "]

res = map(lambda x, y: x.strip() + y.strip(), ls, ls)

print(list(res))


def is_odd(n):
    return n % 2 == 1


tmplist = filter(is_odd, [1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
newlist = list(tmplist)
print(newlist)

c = re.compile(r"[\u4e00-\u9fa5]{2,}")
