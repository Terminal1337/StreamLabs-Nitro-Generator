import json
import re
from typing import Pattern, Union


def escape_regex_flags(pattern: Pattern) -> str:
    flags = ""
    if pattern.flags != 0:
        flags = ""
    if pattern.flags & int(re.IGNORECASE) != 0:
        flags += "i"
    if pattern.flags & int(re.DOTALL) != 0:
        flags += "s"
    if pattern.flags & int(re.MULTILINE) != 0:
        flags += "m"
    assert (
        pattern.flags
        & ~(int(re.MULTILINE) | int(re.IGNORECASE) | int(re.DOTALL) | int(re.UNICODE))
        == 0
    ), "Unexpected re.Pattern flag, only MULTILINE, IGNORECASE and DOTALL are supported."
    return flags


def escape_for_regex(text: str) -> str:
    return re.sub("[.*+?^>${}()|[\\]\\\\]", "\\$&", text)


def escape_regex_for_selector(text: Pattern) -> str:
    return (
        "/"
        + re.sub(
            "(^|[^\\\\])(\\\\\\\\)*([\"\\'`])", "\\1\\2\\\\\\3", text.pattern
        ).replace(">>", "\\>\\>")
        + "/"
        + escape_regex_flags(text)
    )


def escape_for_text_selector(
    text: Union[str, Pattern[str]], exact: bool = None, case_sensitive: bool = None
) -> str:
    if isinstance(text, Pattern):
        return escape_regex_for_selector(text)
    return json.dumps(text) + ("s" if exact else "i")


def escape_for_attribute_selector(
    value: Union[str, Pattern], exact: bool = None
) -> str:
    if isinstance(value, Pattern):
        return escape_regex_for_selector(value)
    return (
        '"'
        + value.replace("\\", "\\\\").replace('"', '\\"')
        + '"'
        + ("s" if exact else "i")
    )
