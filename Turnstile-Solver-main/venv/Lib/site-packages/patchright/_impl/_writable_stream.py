import base64
import os
from pathlib import Path
from typing import Dict, Union
from patchright._impl._connection import ChannelOwner

_WINDOWS = os.name == "nt"
COPY_BUFSIZE = 1024 * 1024 if _WINDOWS else 64 * 1024


class WritableStream(ChannelOwner):

    def __init__(
        self, parent: ChannelOwner, type: str, guid: str, initializer: Dict
    ) -> None:
        super().__init__(parent, type, guid, initializer)

    async def copy(self, path: Union[str, Path]) -> None:
        with open(path, "rb") as f:
            while True:
                data = f.read(COPY_BUFSIZE)
                if not data:
                    break
                await self._channel.send(
                    "write", {"binary": base64.b64encode(data).decode()}
                )
        await self._channel.send("close")
