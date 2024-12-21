import pathlib
from pathlib import Path
from typing import Dict, Optional, Union, cast
from patchright._impl._connection import ChannelOwner, from_channel
from patchright._impl._helper import Error, make_dirs_for_file, patch_error_message
from patchright._impl._stream import Stream


class Artifact(ChannelOwner):

    def __init__(
        self, parent: ChannelOwner, type: str, guid: str, initializer: Dict
    ) -> None:
        super().__init__(parent, type, guid, initializer)
        self.absolute_path = initializer["absolutePath"]

    async def path_after_finished(self) -> pathlib.Path:
        if self._connection.is_remote:
            raise Error(
                "Path is not available when using browser_type.connect(). Use save_as() to save a local copy."
            )
        path = await self._channel.send("pathAfterFinished")
        return pathlib.Path(path)

    async def save_as(self, path: Union[str, Path]) -> None:
        stream = cast(Stream, from_channel(await self._channel.send("saveAsStream")))
        make_dirs_for_file(path)
        await stream.save_as(path)

    async def failure(self) -> Optional[str]:
        reason = await self._channel.send("failure")
        if reason is None:
            return None
        return patch_error_message(reason)

    async def delete(self) -> None:
        await self._channel.send("delete")

    async def read_info_buffer(self) -> bytes:
        stream = cast(Stream, from_channel(await self._channel.send("stream")))
        buffer = await stream.read_all()
        return buffer

    async def cancel(self) -> None:
        await self._channel.send("cancel")
