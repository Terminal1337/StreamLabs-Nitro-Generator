from typing import Dict, Optional
from patchright._impl._connection import Channel
from patchright._impl._element_handle import ElementHandle
from patchright._impl._helper import locals_to_params


def _ax_node_from_protocol(axNode: Dict) -> Dict:
    result = {**axNode}
    if "valueNumber" in axNode:
        result["value"] = axNode["valueNumber"]
    elif "valueString" in axNode:
        result["value"] = axNode["valueString"]
    if "checked" in axNode:
        result["checked"] = (
            True
            if axNode.get("checked") == "checked"
            else (
                False if axNode.get("checked") == "unchecked" else axNode.get("checked")
            )
        )
    if "pressed" in axNode:
        result["pressed"] = (
            True
            if axNode.get("pressed") == "pressed"
            else False if axNode.get("pressed") == "released" else axNode.get("pressed")
        )
    if axNode.get("children"):
        result["children"] = list(map(_ax_node_from_protocol, axNode["children"]))
    if "valueNumber" in result:
        del result["valueNumber"]
    if "valueString" in result:
        del result["valueString"]
    return result


class Accessibility:

    def __init__(self, channel: Channel) -> None:
        self._channel = channel
        self._loop = channel._connection._loop
        self._dispatcher_fiber = channel._connection._dispatcher_fiber

    async def snapshot(
        self, interestingOnly: bool = None, root: ElementHandle = None
    ) -> Optional[Dict]:
        params = locals_to_params(locals())
        if root:
            params["root"] = root._channel
        result = await self._channel.send("accessibilitySnapshot", params)
        return _ax_node_from_protocol(result) if result else None
