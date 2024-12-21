import asyncio
from pathlib import Path
from typing import Any, Dict, List, Set, Union
from patchright._impl._connection import ChannelOwner
from patchright._impl._errors import Error
from patchright._impl._helper import async_readfile
from patchright._impl._locator import set_test_id_attribute_name, test_id_attribute_name


class Selectors:

    def __init__(self, loop: asyncio.AbstractEventLoop, dispatcher_fiber: Any) -> None:
        self._loop = loop
        self._channels: Set[SelectorsOwner] = set()
        self._registrations: List[Dict] = []
        self._dispatcher_fiber = dispatcher_fiber

    async def register(
        self,
        name: str,
        script: str = None,
        path: Union[str, Path] = None,
        contentScript: bool = None,
    ) -> None:
        if not script and (not path):
            raise Error("Either source or path should be specified")
        if path:
            script = (await async_readfile(path)).decode()
        params: Dict[str, Any] = dict(name=name, source=script)
        if contentScript:
            params["contentScript"] = True
        for channel in self._channels:
            await channel._channel.send("register", params)
        self._registrations.append(params)

    def set_test_id_attribute(self, attributeName: str) -> None:
        set_test_id_attribute_name(attributeName)
        for channel in self._channels:
            channel._channel.send_no_reply(
                "setTestIdAttributeName", {"testIdAttributeName": attributeName}
            )

    def _add_channel(self, channel: "SelectorsOwner") -> None:
        self._channels.add(channel)
        for params in self._registrations:
            channel._channel.send_no_reply("register", params)
            channel._channel.send_no_reply(
                "setTestIdAttributeName",
                {"testIdAttributeName": test_id_attribute_name()},
            )

    def _remove_channel(self, channel: "SelectorsOwner") -> None:
        if channel in self._channels:
            self._channels.remove(channel)


class SelectorsOwner(ChannelOwner):

    def __init__(
        self, parent: ChannelOwner, type: str, guid: str, initializer: Dict
    ) -> None:
        super().__init__(parent, type, guid, initializer)
