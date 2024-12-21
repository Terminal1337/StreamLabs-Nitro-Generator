from typing import TYPE_CHECKING, Dict, Optional
from patchright._impl._connection import ChannelOwner, from_nullable_channel
from patchright._impl._helper import locals_to_params

if TYPE_CHECKING:
    from patchright._impl._page import Page


class Dialog(ChannelOwner):

    def __init__(
        self, parent: ChannelOwner, type: str, guid: str, initializer: Dict
    ) -> None:
        super().__init__(parent, type, guid, initializer)
        self._page: Optional["Page"] = from_nullable_channel(initializer.get("page"))

    def __repr__(self) -> str:
        return f"<Dialog type={self.type} message={self.message} default_value={self.default_value}>"

    @property
    def type(self) -> str:
        return self._initializer["type"]

    @property
    def message(self) -> str:
        return self._initializer["message"]

    @property
    def default_value(self) -> str:
        return self._initializer["defaultValue"]

    @property
    def page(self) -> Optional["Page"]:
        return self._page

    async def accept(self, promptText: str = None) -> None:
        await self._channel.send("accept", locals_to_params(locals()))

    async def dismiss(self) -> None:
        await self._channel.send("dismiss")
