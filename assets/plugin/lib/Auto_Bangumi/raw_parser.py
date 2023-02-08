import re

# from parser.episode import Episode

EPISODE_RE = re.compile(r"\d+")
TITLE_RE = re.compile(
    r"(.*|\[.*])( -? \d+|\[\d+]|\[\d+.?[vV]\d{1}]|[第]\d+[话話集]|\[\d+.?END])(.*)"
)
RESOLUTION_RE = re.compile(r"1080|720|2160|4K")
SOURCE_RE = re.compile(r"B-Global|[Bb]aha|[Bb]ilibili|AT-X|Web")
SUB_RE = re.compile(r"[简繁日字幕]|CH|BIG5|GB")

CHINESE_NUMBER_MAP = {
    "一": 1,
    "二": 2,
    "三": 3,
    "四": 4,
    "五": 5,
    "六": 6,
    "七": 7,
    "八": 8,
    "九": 9,
    "十": 10,
}


class Episode:

    def __init__(self, name_en, name_zh, name_jp, season, sr, episode, sub, group, dpi, source):
        self.title_en = name_en
        self.title_zh = name_zh
        self.title_jp = name_jp
        self.season = season
        self.season_raw = sr
        self.episode = episode
        self.sub = sub
        self.group = group
        self.resolution = dpi
        self.source = source

    def __str__(self):
        return "title_en:「%s」, title_zh:「%s」, title_jp:「%s」, season: %s, season_raw: %s, episode: %s, " \
               "sub: %s, resolution: %s, source: %s, group: %s" % \
               (self.title_en, self.title_zh, self.title_jp, self.season, self.season_raw,
                self.episode, self.sub, self.resolution, self.source, self.group)


class RawParser:
    @staticmethod
    def get_group(name):
        sp = re.split(r"[\[\]]", name)
        if len(sp) > 1:
            return sp[1]
        return ""

    @staticmethod
    def pre_process(raw_name):
        return raw_name.replace("【", "[").replace("】", "]")

    @staticmethod
    def season_process(season_info):
        if re.search(r"新番|月?番", season_info):
            name_season = re.sub(".*新番.", "", season_info)
        else:
            name_season = re.sub(r"^[^]】]*[]】]", "", season_info).strip()
        season_rule = r"S\d{1,2}|Season \d{1,2}|[第].[季期]"
        name_season = re.sub(r"[\[\]]", " ", name_season)
        seasons = re.findall(season_rule, name_season)
        if not seasons:
            return name_season, "", 1
        name = re.sub(season_rule, "", name_season)
        for season in seasons:
            season_raw = season
            if re.search(r"S|Season", season) is not None:
                season = int(re.sub(r"S|Season", "", season))
                break
            elif re.search(r"[第 ].*[季期]", season) is not None:
                season_pro = re.sub(r"[第季期 ]", "", season)
                try:
                    season = int(season_pro)
                except ValueError:
                    season = CHINESE_NUMBER_MAP[season_pro]
                    break
        return name, season_raw, season

    @staticmethod
    def name_process(name):
        name_en, name_zh, name_jp = "", "", ""
        name = name.strip()
        split = re.split("/|\s{2}|-\s{2}", name.replace("（仅限港澳台地区）", ""))
        while "" in split:
            split.remove("")
        if len(split) == 1:
            if re.search("_{1}", name) is not None:
                split = re.split("_", name)
            elif re.search(" - {1}", name) is not None:
                split = re.split("-", name)
        if len(split) == 1:
            split_space = name.split(" ")
            for idx, item in enumerate(split_space):
                if re.search("^[\u4e00-\u9fa5]{2,}", item) is not None:
                    split_space.remove(item)
                    split = [item.strip(), " ".join(split_space).strip()]
                    break
        for item in split:
            if re.search("[\u0800-\u4e00]{2,}", item) and not name_jp:
                name_jp = item.strip()
            elif re.search("[\u4e00-\u9fa5]{2,}", item) and not name_zh:
                name_zh = item.strip()
            elif re.search(r"[a-zA-Z]{3,}", item) and not name_en:
                name_en = item.strip()
        return name_en, name_zh, name_jp

    @staticmethod
    def find_tags(other):
        elements = re.sub(r"[\[\]()（）]", " ", other).split(" ")
        # find CHT
        sub, resolution, source = "", "", ""
        for element in filter(lambda x: x != "", elements):
            if SUB_RE.search(element):
                sub = element
            elif RESOLUTION_RE.search(element):
                resolution = element
            elif SOURCE_RE.search(element):
                source = element
        if sub is not None:
            sub = re.sub(r"_MP4|_MKV", "", sub)
        return sub, resolution, source

    def process(self, raw_title):
        raw_title = raw_title.strip()
        content_title = self.pre_process(raw_title)  # 预处理标题
        group = self.get_group(content_title)  # 翻译组的名字
        match_obj = TITLE_RE.match(content_title)  # 处理标题
        season_info, episode_info, other = "", "", ""
        if match_obj:
            season_info, episode_info, other = list(map(
                lambda x: x.strip(), match_obj.groups()
            ))
        raw_name, season_raw, season = self.season_process(season_info)  # 处理 第n季
        name_en, name_zh, name_jp = "", "", ""
        try:
            name_en, name_zh, name_jp = self.name_process(raw_name)  # 处理 名字
        except ValueError:
            pass
        # 处理 集数
        raw_episode = EPISODE_RE.search(episode_info)
        episode = 0
        if raw_episode is not None:
            episode = int(raw_episode.group())
        sub, dpi, source = self.find_tags(other)  # 剩余信息处理
        return name_en, name_zh, name_jp, season, season_raw, episode, sub, dpi, source, group

    def analyse(self, raw):
        try:
            ret = self.process(raw)
            if ret is None:
                return None
            name_en, name_zh, name_jp, season, sr, episode, sub, dpi, source, group = ret
        except Exception as e:
            print('Parser cannot analyse', raw, e)
            return None
        return Episode(name_en, name_zh, name_jp, season, sr, episode, sub, group, dpi, source)


parser = RawParser()


def parse(args):
    ep = parser.analyse(args['title'])
    return ep


def __test__():
    test_txt = "[梦蓝字幕组]New Doraemon 哆啦A梦新番[716][2022.07.23][AVC][10080P][GB_JP]"
    ep = parser.analyse(test_txt)
    print(ep)
    print("en:%s, zh:%s, jp:%s, group:%s" % (ep.title_en, ep.title_zh, ep.title_jp, ep.group))
    return ep
