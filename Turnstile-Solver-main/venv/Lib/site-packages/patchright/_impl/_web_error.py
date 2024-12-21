from asyncio import AbstractEventLoop
from typing import Optional
from patchright._impl._helper import Error
from patchright._impl._page import Page


class WebError:

    def __init__(
        self, loop: AbstractEventLoop, page: Optional[Page], error: Error
    ) -> None:
        self._loop = loop
        self._page = page
        self._error = error

    @property
    def page(self) -> Optional[Page]:
        return self._page

    @property
    def error(self) -> Error:
        return self._error
