# 使用 plugin/lib/Auto_Bangumi/raw_parser.py
import sys
sys.path.append("../../lib/Auto_Bangumi")
import raw_parser


def analyse(item):
    return raw_parser.parser.analyse(item.Name)
