// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package re implements the Python Regular Expression module.
package re

import (
	"github.com/go-python/gpython/py"
)

// Initialise the module
func init() {
	methods := []*py.Method{
		py.MustNewMethod("match", match, 0, re_match_doc),
		py.MustNewMethod("fullmatch", fullmatch, 0, re_fullmatch_doc),
		py.MustNewMethod("search", search, 0, re_search_doc),
		py.MustNewMethod("sub", sub, 0, re_sub_doc),
		py.MustNewMethod("subn", subn, 0, re_subn_doc),
		py.MustNewMethod("split", split, 0, re_split_doc),
		py.MustNewMethod("findall", findall, 0, re_findall_doc),
		py.MustNewMethod("finditer", finditer, 0, re_finditer_doc),
		py.MustNewMethod("compile", compile, 0, re_compile_doc),
		py.MustNewMethod("purge", purge, 0, re_purge_doc),
		py.MustNewMethod("escape", escape, 0, re_escape_doc),
	}
	globals := py.StringDict{
		"I":          py.Int(SRE_FLAG_IGNORECASE),
		"IGNORECASE": py.Int(SRE_FLAG_IGNORECASE),
		"M":          py.Int(SRE_FLAG_MULTILINE),
		"MULTILINE":  py.Int(SRE_FLAG_MULTILINE),
		"S":          py.Int(SRE_FLAG_DOTALL),
		"DOTALL":     py.Int(SRE_FLAG_DOTALL),
	}
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "re",
			Doc:  module_doc,
		},
		Methods: methods,
		Globals: globals,
	})
}

const module_doc = `Support for regular expressions (RE).

This module provides regular expression matching operations similar to
those found in Perl.  It supports both 8-bit and Unicode strings; both
the pattern and the strings being processed can contain null bytes and
characters outside the US ASCII range.

Regular expressions can contain both special and ordinary characters.
Most ordinary characters, like "A", "a", or "0", are the simplest
regular expressions; they simply match themselves.  You can
concatenate ordinary characters, so last matches the string 'last'.

The special characters are:
    "."      Matches any character except a newline.
    "^"      Matches the start of the string.
    "$"      Matches the end of the string or just before the newline at
             the end of the string.
    "*"      Matches 0 or more (greedy) repetitions of the preceding RE.
             Greedy means that it will match as many repetitions as possible.
    "+"      Matches 1 or more (greedy) repetitions of the preceding RE.
    "?"      Matches 0 or 1 (greedy) of the preceding RE.
    *?,+?,?? Non-greedy versions of the previous three special characters.
    {m,n}    Matches from m to n repetitions of the preceding RE.
    {m,n}?   Non-greedy version of the above.
    "\\"     Either escapes special characters or signals a special sequence.
    []       Indicates a set of characters.
             A "^" as the first character indicates a complementing set.
    "|"      A|B, creates an RE that will match either A or B.
    (...)    Matches the RE inside the parentheses.
             The contents can be retrieved or matched later in the string.
    (?aiLmsux) Set the A, I, L, M, S, U, or X flag for the RE (see below).
    (?:...)  Non-grouping version of regular parentheses.
    (?P<name>...) The substring matched by the group is accessible by name.
    (?P=name)     Matches the text matched earlier by the group named name.
    (?#...)  A comment; ignored.
    (?=...)  Matches if ... matches next, but doesn't consume the string.
    (?!...)  Matches if ... doesn't match next.
    (?<=...) Matches if preceded by ... (must be fixed length).
    (?<!...) Matches if not preceded by ... (must be fixed length).
    (?(id/name)yes|no) Matches yes pattern if the group with id/name matched,
                       the (optional) no pattern otherwise.

The special sequences consist of "\\" and a character from the list
below.  If the ordinary character is not on the list, then the
resulting RE will match the second character.
    \number  Matches the contents of the group of the same number.
    \A       Matches only at the start of the string.
    \Z       Matches only at the end of the string.
    \b       Matches the empty string, but only at the start or end of a word.
    \B       Matches the empty string, but not at the start or end of a word.
    \d       Matches any decimal digit; equivalent to the set [0-9] in
             bytes patterns or string patterns with the ASCII flag.
             In string patterns without the ASCII flag, it will match the whole
             range of Unicode digits.
    \D       Matches any non-digit character; equivalent to [^\d].
    \s       Matches any whitespace character; equivalent to [ \t\n\r\f\v] in
             bytes patterns or string patterns with the ASCII flag.
             In string patterns without the ASCII flag, it will match the whole
             range of Unicode whitespace characters.
    \S       Matches any non-whitespace character; equivalent to [^\s].
    \w       Matches any alphanumeric character; equivalent to [a-zA-Z0-9_]
             in bytes patterns or string patterns with the ASCII flag.
             In string patterns without the ASCII flag, it will match the
             range of Unicode alphanumeric characters (letters plus digits
             plus underscore).
             With LOCALE, it will match the set [0-9_] plus characters defined
             as letters for the current locale.
    \W       Matches the complement of \w.
    \\       Matches a literal backslash.

This module exports the following functions:
    match     Match a regular expression pattern to the beginning of a string.
    fullmatch Match a regular expression pattern to all of a string.
    search    Search a string for the presence of a pattern.
    sub       Substitute occurrences of a pattern found in a string.
    subn      Same as sub, but also return the number of substitutions made.
    split     Split a string by the occurrences of a pattern.
    findall   Find all occurrences of a pattern in a string.
    finditer  Return an iterator yielding a match object for each match.
    compile   Compile a pattern into a RegexObject.
    purge     Clear the regular expression cache.
    escape    Backslash all non-alphanumerics in a string.

Some of the functions in this module takes flags as optional parameters:
    A  ASCII       For string patterns, make \w, \W, \b, \B, \d, \D
                   match the corresponding ASCII character categories
                   (rather than the whole Unicode categories, which is the
                   default).
                   For bytes patterns, this flag is the only available
                   behaviour and needn't be specified.
    I  IGNORECASE  Perform case-insensitive matching.
    L  LOCALE      Make \w, \W, \b, \B, dependent on the current locale.
    M  MULTILINE   "^" matches the beginning of lines (after a newline)
                   as well as the string.
                   "$" matches the end of lines (before a newline) as well
                   as the end of the string.
    S  DOTALL      "." matches any character at all, including the newline.
    X  VERBOSE     Ignore whitespace and comments for nicer looking RE's.
    U  UNICODE     For compatibility only. Ignored for string patterns (it
                   is the default), and forbidden for bytes patterns.

This module also defines an exception 'error'.

`

const re_match_doc = `Try to apply the pattern at the start of the string, returning
    a match object, or None if no match was found.`

func match(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var String py.Object
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "string", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|n:match", kwlist, &pattern, &String, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).match(py.Tuple{String}, py.StringDict{})
}

const re_fullmatch_doc = `Try to apply the pattern to all of the string, returning
    a match object, or None if no match was found.`

func fullmatch(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var String py.Object
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "string", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|n:match", kwlist, &pattern, &String, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).fullmatch(py.Tuple{String}, py.StringDict{})
}

const re_search_doc = `Scan through string looking for a match to the pattern, returning
    a match object, or None if no match was found.`

func search(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var String py.Object
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "string", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|n:search", kwlist, &pattern, &String, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).search(py.Tuple{String}, py.StringDict{})
}

const re_sub_doc = `Return the string obtained by replacing the leftmost
    non-overlapping occurrences of the pattern in string by the
    replacement repl.  repl can be either a string or a callable;
    if a string, backslash escapes in it are processed.  If it is
    a callable, it's passed the match object and must return
    a replacement string to be used.`

func sub(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var repl py.Object
	var String py.Object
	var count py.Object = py.Int(0)
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "repl", "string", "count", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OOO|nn:sub", kwlist, &pattern, &repl, &String, &count, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).sub(py.Tuple{repl, String, count}, py.StringDict{})
}

const re_subn_doc = `Return a 2-tuple containing (new_string, number).
    new_string is the string obtained by replacing the leftmost
    non-overlapping occurrences of the pattern in the source
    string by the replacement repl.  number is the number of
    substitutions that were made. repl can be either a string or a
    callable; if a string, backslash escapes in it are processed.
    If it is a callable, it's passed the match object and must
    return a replacement string to be used.`

func subn(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var repl py.Object
	var String py.Object
	var count py.Object = py.Int(0)
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "repl", "string", "count", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OOO|nn:sub", kwlist, &pattern, &repl, &String, &count, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).subn(py.Tuple{repl, String, count}, py.StringDict{})
}

const re_split_doc = `Split the source string by the occurrences of the pattern,
    returning a list containing the resulting substrings.  If
    capturing parentheses are used in pattern, then the text of all
    groups in the pattern are also returned as part of the resulting
    list.  If maxsplit is nonzero, at most maxsplit splits occur,
    and the remainder of the string is returned as the final element
    of the list.`

func split(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var String py.Object
	var maxSplit py.Object = py.Int(0)
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "string", "maxsplit", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|nn:sub", kwlist, &pattern, &String, &maxSplit, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).split(py.Tuple{String, maxSplit}, py.StringDict{})
}

const re_findall_doc = `Return a list of all non-overlapping matches in the string.

    If one or more capturing groups are present in the pattern, return
    a list of groups; this will be a list of tuples if the pattern
    has more than one group.

    Empty matches are included in the result.`

func findall(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var String py.Object
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "string", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|n:findall", kwlist, &pattern, &String, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).findall(py.Tuple{String}, py.StringDict{})
}

const re_finditer_doc = `Return an iterator over all non-overlapping matches in the
        string.  For each match, the iterator returns a match object.

        Empty matches are included in the result.`

func finditer(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var String py.Object
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "string", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|n:findall", kwlist, &pattern, &String, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags).finditer(py.Tuple{String}, py.StringDict{})
}

const re_compile_doc = `Compile a regular expression pattern, returning a pattern object.`

func compile(self py.Object, args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var pattern py.Object
	var flags py.Object = py.Int(0)
	kwlist := []string{"pattern", "flags"}
	err = py.ParseTupleAndKeywords(args, kwargs, "O|n:findall", kwlist, &pattern, &flags)
	if err != nil {
		return nil, err
	}
	return PatternNew(pattern, flags), nil
}

const re_purge_doc = `Clear the regular expression caches`

func purge(self py.Object) (py.Object, error) {
	//_cache.clear()
	//_cache_repl.clear()
	return nil, nil
}

const re_escape_doc = `Escape all the characters in pattern except ASCII letters, numbers and '_'.`

func escape(self py.Object, args py.Object) (py.Object, error) {
	str, _ := toString(args)
	return py.String(py.StringEscape(py.String(str), true)), nil

}
