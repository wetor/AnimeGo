import re
import libtest as self

S = str
B = bytes

doc="test_search_star_plus"

self.assertEqual(re.search('x*', 'axx').span(0), (0, 0))
self.assertEqual(re.search('x*', 'axx').span(), (0, 0))
self.assertEqual(re.search('x+', 'axx').span(0), (1, 3))
self.assertEqual(re.search('x+', 'axx').span(), (1, 3))
self.assertIsNone(re.search('x', 'aaa'))
self.assertEqual(re.match('a*', 'xxx').span(0), (0, 0))
self.assertEqual(re.match('a*', 'xxx').span(), (0, 0))
self.assertEqual(re.match('x*', 'xxxa').span(0), (0, 3))
self.assertEqual(re.match('x*', 'xxxa').span(), (0, 3))
self.assertIsNone(re.match('a+', 'xxx'))

doc="test_basic_re_sub"
def bump_num(matchobj):
    int_value = int(matchobj.group(0))
    return str(int_value + 1)

self.assertTypedEqual(re.sub('y', 'a', 'xyz'), 'xaz')
self.assertTypedEqual(re.sub('y', S('a'), S('xyz')), 'xaz')
self.assertTypedEqual(re.sub(b'y', b'a', b'xyz'), b'xaz')
self.assertTypedEqual(re.sub(b'y', B(b'a'), B(b'xyz')), b'xaz')
# FIXME self.assertTypedEqual(re.sub(b'y', bytearray(b'a'), bytearray(b'xyz')), b'xaz')
# FIXME self.assertTypedEqual(re.sub(b'y', memoryview(b'a'), memoryview(b'xyz')), b'xaz')
for y in ("\xe0", "\u0430", "\U0001d49c"):
    self.assertEqual(re.sub(y, 'a', 'x%sz' % y), 'xaz')

self.assertEqual(re.sub("(?i)b+", "x", "bbbb BBBB"), 'x x')
self.assertEqual(re.sub(r'\d+', bump_num, '08.2 -2 23x99y'),
                 '9.3 -3 24x100y')
self.assertEqual(re.sub(r'\d+', bump_num, '08.2 -2 23x99y', 3),
                 '9.3 -3 23x99y')

self.assertEqual(re.sub('.', lambda m: r"\n", 'x'), '\\n')
self.assertEqual(re.sub('.', r"\n", 'x'), '\n')

s = r"\1\1"
# FIXME self.assertEqual(re.sub('(.)', s, 'x'), 'xx')
# FIXME self.assertEqual(re.sub('(.)', re.escape(s), 'x'), s)
self.assertEqual(re.sub('(.)', lambda m: s, 'x'), s)

# FIXME self.assertEqual(re.sub('(?P<a>x)', '\g<a>\g<a>', 'xx'), 'xxxx')
# FIXME self.assertEqual(re.sub('(?P<a>x)', '\g<a>\g<1>', 'xx'), 'xxxx')
# FIXME self.assertEqual(re.sub('(?P<unk>x)', '\g<unk>\g<unk>', 'xx'), 'xxxx')
# FIXME self.assertEqual(re.sub('(?P<unk>x)', '\g<1>\g<1>', 'xx'), 'xxxx')

self.assertEqual(re.sub('a',r'\t\n\v\r\f\a\b\B\Z\a\A\w\W\s\S\d\D','a'),
                 '\t\n\v\r\f\a\b\\B\\Z\a\\A\\w\\W\\s\\S\\d\\D')
self.assertEqual(re.sub('a', '\t\n\v\r\f\a', 'a'), '\t\n\v\r\f\a')
self.assertEqual(re.sub('a', '\t\n\v\r\f\a', 'a'),
                 (chr(9)+chr(10)+chr(11)+chr(13)+chr(12)+chr(7)))

self.assertEqual(re.sub('^\s*', 'X', 'test'), 'Xtest')

doc="test_bug_449964"
# fails for group followed by other escape
# FIXME self.assertEqual(re.sub(r'(?P<unk>x)', '\g<1>\g<1>\\b', 'xx'),
# FIXME                  'xx\bxx\b')
doc="test_bug_449000"
# Test for sub() on escaped characters
self.assertEqual(re.sub(r'\r\n', r'\n', 'abc\r\ndef\r\n'),
                 'abc\ndef\n')
self.assertEqual(re.sub('\r\n', r'\n', 'abc\r\ndef\r\n'),
                 'abc\ndef\n')
self.assertEqual(re.sub(r'\r\n', '\n', 'abc\r\ndef\r\n'),
                 'abc\ndef\n')
self.assertEqual(re.sub('\r\n', '\n', 'abc\r\ndef\r\n'),
                 'abc\ndef\n')

doc="test_bug_1661"
# Verify that flags do not get silently ignored with compiled patterns
pattern = re.compile('.')
# FIXME self.assertRaises(ValueError, re.match, pattern, 'A', re.I)
# FIXME self.assertRaises(ValueError, re.search, pattern, 'A', re.I)
# FIXME self.assertRaises(ValueError, re.findall, pattern, 'A', re.I)
# FIXME self.assertRaises(ValueError, re.compile, pattern, re.I)
doc="test_bug_3629"
# A regex that triggered a bug in the sre-code validator
# FIXME re.compile("(?P<quote>)(?(quote))")
doc="test_sub_template_numeric_escape"
doc="test_qualified_re_sub"
self.assertEqual(re.sub('a', 'b', 'aaaaa'), 'bbbbb')
self.assertEqual(re.sub('a', 'b', 'aaaaa', 1), 'baaaa')
doc="test_bug_114660"
doc="test_bug_462270"
# Test for empty sub() behaviour, see SF bug #462270
self.assertEqual(re.sub('x*', '-', 'abxd'), '-a-b-d-')
self.assertEqual(re.sub('x+', '-', 'abxd'), 'ab-d')
doc="test_symbolic_groups"
doc="test_symbolic_refs"
doc="test_re_subn"
self.assertEqual(re.subn("(?i)b+", "x", "bbbb BBBB"), ('x x', 2))
self.assertEqual(re.subn("b+", "x", "bbbb BBBB"), ('x BBBB', 1))
self.assertEqual(re.subn("b+", "x", "xyz"), ('xyz', 0))
self.assertEqual(re.subn("b*", "x", "xyz"), ('xxxyxzx', 4))
self.assertEqual(re.subn("b*", "x", "xyz", 2), ('xxxyz', 2))
doc="test_re_split"
for string in ":a:b::c", S(":a:b::c"):
    self.assertTypedEqual(re.split(":", string),
                          ['', 'a', 'b', '', 'c'])
    self.assertTypedEqual(re.split(":*", string),
                          ['', 'a', 'b', 'c'])
    # FIXME self.assertTypedEqual(re.split("(:*)", string),
    # FIXME                       ['', ':', 'a', ':', 'b', '::', 'c'])
for string in (b":a:b::c", B(b":a:b::c")):
    self.assertTypedEqual(re.split(b":", string),
                          [b'', b'a', b'b', b'', b'c'])
    self.assertTypedEqual(re.split(b":*", string),
                          [b'', b'a', b'b', b'c'])
    # FIXME self.assertTypedEqual(re.split(b"(:*)", string),
    # FIXME                       [b'', b':', b'a', b':', b'b', b'::', b'c'])
for a, b, c in ("\xe0\xdf\xe7", "\u0430\u0431\u0432",
                "\U0001d49c\U0001d49e\U0001d4b5"):
    string = ":%s:%s::%s" % (a, b, c)
    self.assertEqual(re.split(":", string), ['', a, b, '', c])
    self.assertEqual(re.split(":*", string), ['', a, b, c])
    # FIXME self.assertEqual(re.split("(:*)", string),
    # FIXME                  ['', ':', a, ':', b, '::', c])

self.assertEqual(re.split("(?::*)", ":a:b::c"), ['', 'a', 'b', 'c'])
# FIXME self.assertEqual(re.split("(:)*", ":a:b::c"),
# FIXME                  ['', ':', 'a', ':', 'b', ':', 'c'])
# FIXME self.assertEqual(re.split("([b:]+)", ":a:b::c"),
# FIXME                  ['', ':', 'a', ':b::', 'c'])
# FIXME self.assertEqual(re.split("(b)|(:+)", ":a:b::c"),
# FIXME                  ['', None, ':', 'a', None, ':', '', 'b', None, '',
# FIXME                   None, '::', 'c'])
self.assertEqual(re.split("(?:b)|(?::+)", ":a:b::c"),
                 ['', 'a', '', '', 'c'])
doc="test_qualified_re_split"
self.assertEqual(re.split(":", ":a:b::c", 2), ['', 'a', 'b::c'])
self.assertEqual(re.split(':', 'a:b:c:d', 2), ['a', 'b', 'c:d'])
# FIXME self.assertEqual(re.split("(:)", ":a:b::c", 2),
# FIXME                  ['', ':', 'a', ':', 'b::c'])
# FIXME self.assertEqual(re.split("(:*)", ":a:b::c", 2),
# FIXME                  ['', ':', 'a', ':', 'b::c'])
doc="test_re_findall"
self.assertEqual(re.findall(":+", "abc"), [])
for string in "a:b::c:::d", S("a:b::c:::d"):
    self.assertTypedEqual(re.findall(":+", string),
                          [":", "::", ":::"])
    # FIXME self.assertTypedEqual(re.findall("(:+)", string),
    # FIXME                       [":", "::", ":::"])
    # FIXME self.assertTypedEqual(re.findall("(:)(:*)", string),
    # FIXME                       [(":", ""), (":", ":"), (":", "::")])
for string in (b"a:b::c:::d", B(b"a:b::c:::d")):
    self.assertTypedEqual(re.findall(b":+", string),
                          [b":", b"::", b":::"])
    # FIXME self.assertTypedEqual(re.findall(b"(:+)", string),
    # FIXME                       [b":", b"::", b":::"])
    # FIXME self.assertTypedEqual(re.findall(b"(:)(:*)", string),
    # FIXME                       [(b":", b""), (b":", b":"), (b":", b"::")])
for x in ("\xe0", "\u0430", "\U0001d49c"):
    xx = x * 2
    xxx = x * 3
    string = "a%sb%sc%sd" % (x, xx, xxx)
    self.assertEqual(re.findall("%s+" % x, string), [x, xx, xxx])
    # FIXME self.assertEqual(re.findall("(%s+)" % x, string), [x, xx, xxx])
    # FIXME self.assertEqual(re.findall("(%s)(%s*)" % (x, x), string),
    # FIXME                  [(x, ""), (x, x), (x, xx)])

doc="test_bug_117612"
self.assertEqual(re.findall(r"(a|(b))", "aba"),
                 [("a", ""),("b", "b"),("a", "")])
doc="test_re_match"
for string in 'a', S('a'):
    self.assertEqual(re.match('a', string).groups(), ())
    self.assertEqual(re.match('(a)', string).groups(), ('a',))
    self.assertEqual(re.match('(a)', string).group(0), 'a')
    self.assertEqual(re.match('(a)', string).group(1), 'a')
    self.assertEqual(re.match('(a)', string).group(1, 1), ('a', 'a'))
for string in b'a', B(b'a'):
    self.assertEqual(re.match(b'a', string).groups(), ())
    self.assertEqual(re.match(b'(a)', string).groups(), (b'a',))
    self.assertEqual(re.match(b'(a)', string).group(0), b'a')
    self.assertEqual(re.match(b'(a)', string).group(1), b'a')
    self.assertEqual(re.match(b'(a)', string).group(1, 1), (b'a', b'a'))
for a in ("\xe0", "\u0430", "\U0001d49c"):
    self.assertEqual(re.match(a, a).groups(), ())
    self.assertEqual(re.match('(%s)' % a, a).groups(), (a,))
    self.assertEqual(re.match('(%s)' % a, a).group(0), a)
    self.assertEqual(re.match('(%s)' % a, a).group(1), a)
    self.assertEqual(re.match('(%s)' % a, a).group(1, 1), (a, a))

pat = re.compile('((a)|(b))(c)?')

self.assertEqual(pat.match('a').groups(), ('a', 'a', None, None))
self.assertEqual(pat.match('b').groups(), ('b', None, 'b', None))
self.assertEqual(pat.match('ac').groups(), ('a', 'a', None, 'c'))
self.assertEqual(pat.match('bc').groups(), ('b', None, 'b', 'c'))
self.assertEqual(pat.match('bc').groups(""), ('b', "", 'b', 'c'))

# A single group
m = re.match('(a)', 'a')
self.assertEqual(m.group(0), 'a')
self.assertEqual(m.group(0), 'a')
self.assertEqual(m.group(1), 'a')
self.assertEqual(m.group(1, 1), ('a', 'a'))

pat = re.compile('(?:(?P<a1>a)|(?P<b2>b))(?P<c3>c)?')
self.assertEqual(pat.match('a').group(1, 2, 3), ('a', None, None))
self.assertEqual(pat.match('b').group('a1', 'b2', 'c3'),
                 (None, 'b', None))
self.assertEqual(pat.match('ac').group(1, 'b2', 3), ('a', None, 'c'))

doc="test_re_fullmatch"
# Issue 16203: Proposal: add re.fullmatch() method.
self.assertEqual(re.fullmatch(r"a", "a").span(), (0, 1))
for string in "ab", S("ab"):
    self.assertEqual(re.fullmatch(r"a|ab", string).span(), (0, 2))
for string in b"ab", B(b"ab"):
    self.assertEqual(re.fullmatch(br"a|ab", string).span(), (0, 2))
for a, b in "\xe0\xdf", "\u0430\u0431", "\U0001d49c\U0001d49e":
    r = r"%s|%s" % (a, a + b)
    # FIXME self.assertEqual(re.fullmatch(r, a + b).span(), (0, 2))
self.assertEqual(re.fullmatch(r".*?$", "abc").span(), (0, 3))
self.assertEqual(re.fullmatch(r".*?", "abc").span(), (0, 3))
self.assertEqual(re.fullmatch(r"a.*?b", "ab").span(), (0, 2))
self.assertEqual(re.fullmatch(r"a.*?b", "abb").span(), (0, 3))
self.assertEqual(re.fullmatch(r"a.*?b", "axxb").span(), (0, 4))
# FIXME self.assertIsNone(re.fullmatch(r"a+", "ab"))
self.assertIsNone(re.fullmatch(r"abc$", "abc\n"))
# FIXME self.assertIsNone(re.fullmatch(r"abc\Z", "abc\n"))
# FIXME self.assertIsNone(re.fullmatch(r"(?m)abc$", "abc\n"))
# FIXME self.assertEqual(re.fullmatch(r"ab(?=c)cd", "abcd").span(), (0, 4))
# FIXME self.assertEqual(re.fullmatch(r"ab(?<=b)cd", "abcd").span(), (0, 4))
# FIXME self.assertEqual(re.fullmatch(r"(?=a|ab)ab", "ab").span(), (0, 2))
self.assertEqual(
    re.compile(r"bc").fullmatch("abcd", pos=1, endpos=3).span(), (1, 3))
self.assertEqual(
    re.compile(r".*?$").fullmatch("abcd", pos=1, endpos=3).span(), (1, 3))
self.assertEqual(
    re.compile(r".*?").fullmatch("abcd", pos=1, endpos=3).span(), (1, 3))

doc="test_re_groupref_exists"
doc="test_re_groupref"
doc="test_groupdict"
doc="test_expand"
doc="test_repeat_minmax"
self.assertIsNone(re.match("^(\w){1}$", "abc"))
self.assertIsNone(re.match("^(\w){1}?$", "abc"))
self.assertIsNone(re.match("^(\w){1,2}$", "abc"))
self.assertIsNone(re.match("^(\w){1,2}?$", "abc"))

self.assertEqual(re.match("^(\w){3}$", "abc").group(1), "c")
self.assertEqual(re.match("^(\w){1,3}$", "abc").group(1), "c")
self.assertEqual(re.match("^(\w){1,4}$", "abc").group(1), "c")
self.assertEqual(re.match("^(\w){3,4}?$", "abc").group(1), "c")
self.assertEqual(re.match("^(\w){3}?$", "abc").group(1), "c")
self.assertEqual(re.match("^(\w){1,3}?$", "abc").group(1), "c")
self.assertEqual(re.match("^(\w){1,4}?$", "abc").group(1), "c")
self.assertEqual(re.match("^(\w){3,4}?$", "abc").group(1), "c")

self.assertIsNone(re.match("^x{1}$", "xxx"))
self.assertIsNone(re.match("^x{1}?$", "xxx"))
self.assertIsNone(re.match("^x{1,2}$", "xxx"))
self.assertIsNone(re.match("^x{1,2}?$", "xxx"))

self.assertTrue(re.match("^x{3}$", "xxx"))
self.assertTrue(re.match("^x{1,3}$", "xxx"))
self.assertTrue(re.match("^x{1,4}$", "xxx"))
self.assertTrue(re.match("^x{3,4}?$", "xxx"))
self.assertTrue(re.match("^x{3}?$", "xxx"))
self.assertTrue(re.match("^x{1,3}?$", "xxx"))
self.assertTrue(re.match("^x{1,4}?$", "xxx"))
self.assertTrue(re.match("^x{3,4}?$", "xxx"))

self.assertIsNone(re.match("^x{}$", "xxx"))
self.assertTrue(re.match("^x{}$", "x{}"))

doc="test_getattr"
# FIXME self.assertEqual(re.compile("(?i)(a)(b)").pattern, "(?i)(a)(b)")
# FIXME self.assertEqual(re.compile("(?i)(a)(b)").flags, re.I | re.U)
# FIXME self.assertEqual(re.compile("(?i)(a)(b)").groups, 2)
# FIXME self.assertEqual(re.compile("(?i)(a)(b)").groupindex, {})
# FIXME self.assertEqual(re.compile("(?i)(?P<first>a)(?P<other>b)").groupindex,
# FIXME                  {'first': 1, 'other': 2})

# FIXME self.assertEqual(re.match("(a)", "a").pos, 0)
# FIXME self.assertEqual(re.match("(a)", "a").endpos, 1)
# FIXME self.assertEqual(re.match("(a)", "a").string, "a")
# FIXME self.assertEqual(re.match("(a)", "a").regs, ((0, 1), (0, 1)))
# FIXME self.assertTrue(re.match("(a)", "a").re)
doc="test_special_escapes"
doc="test_string_boundaries"
doc="test_bigcharset"
self.assertEqual(re.match("([\u2222\u2223])",
                          "\u2222").group(1), "\u2222")
# FIXME r = '[%s]' % ''.join(map(chr, range(256, 2**16, 255)))
# FIXME self.assertEqual(re.match(r, "\uff01").group(), "\uff01")
doc="test_big_codesize"
# Issue #1160
r = re.compile('|'.join(list(('%d'%x for x in range(10000)))))
self.assertTrue(r.match('1000'))
self.assertTrue(r.match('9999'))
doc="test_anyall"
self.assertEqual(re.match("a.b", "a\nb", re.DOTALL).group(0),
                 "a\nb")
self.assertEqual(re.match("a.*b", "a\n\nb", re.DOTALL).group(0),
                 "a\n\nb")
doc="test_lookahead"
doc="test_lookbehind"
doc="test_ignore_case"
self.assertEqual(re.match("abc", "ABC", re.I).group(0), "ABC")
self.assertEqual(re.match(b"abc", b"ABC", re.I).group(0), b"ABC")
self.assertEqual(re.match(r"(a\s[^a])", "a b", re.I).group(1), "a b")
self.assertEqual(re.match(r"(a\s[^a]*)", "a bb", re.I).group(1), "a bb")
self.assertEqual(re.match(r"(a\s[abc])", "a b", re.I).group(1), "a b")
self.assertEqual(re.match(r"(a\s[abc]*)", "a bb", re.I).group(1), "a bb")
# FIXME self.assertEqual(re.match(r"((a)\s\2)", "a a", re.I).group(1), "a a")
# FIXME self.assertEqual(re.match(r"((a)\s\2*)", "a aa", re.I).group(1), "a aa")
self.assertEqual(re.match(r"((a)\s(abc|a))", "a a", re.I).group(1), "a a")
self.assertEqual(re.match(r"((a)\s(abc|a)*)", "a aa", re.I).group(1), "a aa")

# FIXME assert '\u212a'.lower() == 'k' # 'K'
self.assertTrue(re.match(r'K', '\u212a', re.I))
self.assertTrue(re.match(r'k', '\u212a', re.I))
# FIXME self.assertTrue(re.match(r'\u212a', 'K', re.I))
# FIXME self.assertTrue(re.match(r'\u212a', 'k', re.I))
# FIXME assert '\u017f'.upper() == 'S' # 'ſ'
self.assertTrue(re.match(r'S', '\u017f', re.I))
self.assertTrue(re.match(r's', '\u017f', re.I))
# FIXME self.assertTrue(re.match(r'\u017f', 'S', re.I))
# FIXME self.assertTrue(re.match(r'\u017f', 's', re.I))
# FIXME assert '\ufb05'.upper() == '\ufb06'.upper() == 'ST' # 'ﬅ', 'ﬆ'
# FIXME self.assertTrue(re.match(r'\ufb05', '\ufb06', re.I))
# FIXME self.assertTrue(re.match(r'\ufb06', '\ufb05', re.I))
doc="test_ignore_case_set"
self.assertTrue(re.match(r'[19A]', 'A', re.I))
self.assertTrue(re.match(r'[19a]', 'a', re.I))
self.assertTrue(re.match(r'[19a]', 'A', re.I))
self.assertTrue(re.match(r'[19A]', 'a', re.I))
self.assertTrue(re.match(br'[19A]', b'A', re.I))
self.assertTrue(re.match(br'[19a]', b'a', re.I))
self.assertTrue(re.match(br'[19a]', b'A', re.I))
self.assertTrue(re.match(br'[19A]', b'a', re.I))
# FIXME assert '\u212a'.lower() == 'k' # 'K'
self.assertTrue(re.match(r'[19K]', '\u212a', re.I))
self.assertTrue(re.match(r'[19k]', '\u212a', re.I))
# FIXME self.assertTrue(re.match(r'[19\u212a]', 'K', re.I))
# FIXME self.assertTrue(re.match(r'[19\u212a]', 'k', re.I))
# FIXME assert '\u017f'.upper() == 'S' # 'ſ'
self.assertTrue(re.match(r'[19S]', '\u017f', re.I))
self.assertTrue(re.match(r'[19s]', '\u017f', re.I))
# FIXME self.assertTrue(re.match(r'[19\u017f]', 'S', re.I))
# FIXME self.assertTrue(re.match(r'[19\u017f]', 's', re.I))
# FIXME assert '\ufb05'.upper() == '\ufb06'.upper() == 'ST' # 'ﬅ', 'ﬆ'
# FIXME self.assertTrue(re.match(r'[19\ufb05]', '\ufb06', re.I))
# FIXME self.assertTrue(re.match(r'[19\ufb06]', '\ufb05', re.I))
doc="test_category"
self.assertEqual(re.match(r"(\s)", " ").group(1), " ")
doc="test_getlower"
# FIXME import _sre
# FIXME self.assertEqual(_sre.getlower(ord('A'), 0), ord('a'))
# FIXME self.assertEqual(_sre.getlower(ord('A'), re.LOCALE), ord('a'))
# FIXME self.assertEqual(_sre.getlower(ord('A'), re.UNICODE), ord('a'))

self.assertEqual(re.match("abc", "ABC", re.I).group(0), "ABC")
self.assertEqual(re.match(b"abc", b"ABC", re.I).group(0), b"ABC")
doc="test_not_literal"
self.assertEqual(re.search("\s([^a])", " b").group(1), "b")
self.assertEqual(re.search("\s([^a]*)", " bb").group(1), "bb")
doc="test_search_coverage"
self.assertEqual(re.search("\s(b)", " b").group(1), "b")
self.assertEqual(re.search("a\s", "a ").group(0), "a ")


def assertMatch(self, pattern, text, match=None, span=None,
                matcher=re.match):
    if match is None and span is None:
        # the pattern matches the whole text
        match = text
        span = (0, len(text))
    elif match is None or span is None:
        raise ValueError('If match is not None, span should be specified '
                         '(and vice versa).')
    m = matcher(pattern, text)
    self.assertTrue(m)
    self.assertEqual(m.group(), match)
    self.assertEqual(m.span(), span)

doc="test_re_escape"
doc="test_re_escape_byte"
doc="test_re_escape_non_ascii"
doc="test_re_escape_non_ascii_bytes"
doc="test_pickling"
doc="test_constants"
doc="test_flags"
doc="test_sre_character_literals"
doc="test_sre_character_class_literals"
doc="test_sre_byte_literals"
doc="test_sre_byte_class_literals"
doc="test_bug_113254"
self.assertEqual(re.match(r'(a)|(b)', 'b').start(1), -1)
self.assertEqual(re.match(r'(a)|(b)', 'b').end(1), -1)
self.assertEqual(re.match(r'(a)|(b)', 'b').span(1), (-1, -1))
doc="test_bug_527371"
doc="test_bug_545855"
dpc="test_bug_418626"
# bugs 418626 at al. -- Testing Greg Chapman's addition of op code
# SRE_OP_MIN_REPEAT_ONE for eliminating recursion on simple uses of
# pattern '*?' on a long string.
self.assertEqual(re.match('.*?c', 10000*'ab'+'cd').end(0), 20001)
self.assertEqual(re.match('.*?cd', 5000*'ab'+'c'+5000*'ab'+'cde').end(0),
                 20003)
self.assertEqual(re.match('.*?cd', 20000*'abc'+'de').end(0), 60001)
# non-simple '*?' still used to hit the recursion limit, before the
# non-recursive scheme was implemented.
self.assertEqual(re.search('(a|b)*?c', 10000*'ab'+'cd').end(0), 20001)
doc="test_bug_612074"
doc="test_stack_overflow"
# nasty cases that used to overflow the straightforward recursive
# implementation of repeated groups.
self.assertEqual(re.match('(x)*', 50000*'x').group(1), 'x')
self.assertEqual(re.match('(x)*y', 50000*'x'+'y').group(1), 'x')
self.assertEqual(re.match('(x)*?y', 50000*'x'+'y').group(1), 'x')
doc="test_unlimited_zero_width_repeat"
# Issue #9669
self.assertIsNone(re.match(r'(?:a?)*y', 'z'))
self.assertIsNone(re.match(r'(?:a?)+y', 'z'))
self.assertIsNone(re.match(r'(?:a?){2,}y', 'z'))
self.assertIsNone(re.match(r'(?:a?)*?y', 'z'))
self.assertIsNone(re.match(r'(?:a?)+?y', 'z'))
self.assertIsNone(re.match(r'(?:a?){2,}?y', 'z'))
doc="test_scanner"
doc="test_bug_448951"
# bug 448951 (similar to 429357, but with single char match)
# (Also test greedy matches.)
for op in '','?','*':
    self.assertEqual(re.match(r'((.%s):)?z'%op, 'z').groups(),
                     (None, None))
    self.assertEqual(re.match(r'((.%s):)?z'%op, 'a:z').groups(),
                     ('a:', 'a'))
dpc="test_bug_725106"
# capturing groups in alternatives in repeats
self.assertEqual(re.match('^((a)|b)*', 'abc').groups(),
                 ('b', 'a'))
self.assertEqual(re.match('^(([ab])|c)*', 'abc').groups(),
                 ('c', 'b'))
self.assertEqual(re.match('^((d)|[ab])*', 'abc').groups(),
                 ('b', None))
self.assertEqual(re.match('^((a)c|[ab])*', 'abc').groups(),
                 ('b', None))
self.assertEqual(re.match('^((a)|b)*?c', 'abc').groups(),
                 ('b', 'a'))
self.assertEqual(re.match('^(([ab])|c)*?d', 'abcd').groups(),
                 ('c', 'b'))
self.assertEqual(re.match('^((d)|[ab])*?c', 'abc').groups(),
                 ('b', None))
self.assertEqual(re.match('^((a)c|[ab])*?c', 'abc').groups(),
                 ('b', None))
doc="test_bug_725149"
doc="test_bug_764548"
doc="test_finditer"
doc="test_bug_926075"
doc="test_bug_931848"
doc="test_bug_581080"
doc="test_bug_817234"
doc="test_bug_6561"
doc="test_empty_array"
doc="test_inline_flags"
doc="test_dollar_matches_twice"
doc="test_bytes_str_mixing"
doc="test_ascii_and_unicode_flag"


doc="finished"
