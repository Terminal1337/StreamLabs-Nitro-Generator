import asyncio
from typing import Any, Generic, TypeVar

T = TypeVar("T")


class EventContextManagerImpl(Generic[T]):

    def __init__(self, future: asyncio.Future) -> None:
        self._future: asyncio.Future = future

    @property
    def future(self) -> asyncio.Future:
        return self._future

    async def __aenter__(self) -> asyncio.Future:
        return self._future

    async def __aexit__(self, *args: Any) -> None:
        await self._future
