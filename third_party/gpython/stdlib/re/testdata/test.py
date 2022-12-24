# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import re

EPISODE_RE = re.compile(r"\d+")
TITLE_RE = re.compile(
    r"(.*|\[.*])( -? \d+ |\[\d+]|\[\d+.?[vV]\d{1}]|[第]\d+[话話集]|\[\d+.?END])(.*)"
)
RESOLUTION_RE = re.compile(r"1080|720|2160|4K")
SOURCE_RE = re.compile(r"B-Global|[Bb]aha|[Bb]ilibili|AT-X|Web")
SUB_RE = re.compile(r"[简繁日字幕]|CH|BIG5|GB")

txt = "[梦蓝字幕组]New Doraemon 哆啦A梦新番[716][2022.07.23][AVC][10080P][GB_JP]"

# match_obj = TITLE_RE.match(txt)  # 处理标题
#
# for group in match_obj.groups():
#     print(group)

print(re.split(r"[一-龥]{2,}", "asd梦蓝字幕组aasd"))


