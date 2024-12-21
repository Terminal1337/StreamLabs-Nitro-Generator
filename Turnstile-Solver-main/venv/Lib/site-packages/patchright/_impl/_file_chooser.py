from pathlib import Path
from typing import TYPE_CHECKING, Sequence, Union
from patchright._impl._api_structures import FilePayload

if TYPE_CHECKING:
    from patchright._impl._element_handle import ElementHandle
    from patchright._impl._page import Page


class FileChooser:

    def __init__(
        self, page: "Page", element_handle: "ElementHandle", is_multiple: bool
    ) -> None:
        self._page = page
        self._loop = page._loop
        self._dispatcher_fiber = page._dispatcher_fiber
        self._element_handle = element_handle
        self._is_multiple = is_multiple

    def __repr__(self) -> str:
        return f"<FileChooser page={self._page} element={self._element_handle}>"

    @property
    def page(self) -> "Page":
        return self._page

    @property
    def element(self) -> "ElementHandle":
        return self._element_handle

    def is_multiple(self) -> bool:
        return self._is_multiple

    async def set_files(
        self,
        files: Union[
            str, Path, FilePayload, Sequence[Union[str, Path]], Sequence[FilePayload]
        ],
        timeout: float = None,
        noWaitAfter: bool = None,
    ) -> None:
        await self._element_handle.set_input_files(files, timeout, noWaitAfter)
