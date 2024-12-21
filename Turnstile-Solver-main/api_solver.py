import asyncio
import time
from typing import Dict, Optional
from dataclasses import dataclass
from patchright.async_api import async_playwright, Page, BrowserContext
from logmagix import Logger, Loader
from quart import Quart, request, jsonify
from collections import deque

@dataclass
class TurnstileResult:
    turnstile_value: Optional[str]
    elapsed_time_seconds: float
    status: str
    reason: Optional[str] = None

@dataclass
class TurnstileAPIResult:
    result: Optional[str]
    elapsed_time_seconds: Optional[float] = None
    status: str = "success"
    error: Optional[str] = None


class PagePool:
    def __init__(self, context, debug: bool = False, log=None):
        self.context = context
        self.min_size = 1
        self.max_size = 10
        self.available_pages: deque = deque()
        self.in_use_pages: set = set()
        self._lock = asyncio.Lock()
        self.debug = debug
        self.log = log

    async def initialize(self):
        """Create initial pool of pages"""
        for _ in range(self.min_size):
            page = await self.context.new_page()
            self.available_pages.append(page)

    async def get_page(self):
        """Get an available page from the pool or create a new one if needed"""
        async with self._lock:
            if (not self.available_pages and
                    len(self.in_use_pages) < self.max_size):
                page = await self.context.new_page()
                if self.debug:
                    self.log.debug(f"Created new page. Total pages: {len(self.in_use_pages) + 1}")
            else:
                while not self.available_pages:
                    await asyncio.sleep(0.1)
                page = self.available_pages.popleft()

            self.in_use_pages.add(page)
            return page

    async def return_page(self, page):
        """Return a page to the pool or close it if we have too many"""
        async with self._lock:
            self.in_use_pages.remove(page)
            total_pages = len(self.available_pages) + len(self.in_use_pages) + 1
            if total_pages > self.min_size and len(self.available_pages) >= 2:
                await page.close()
                if self.debug:
                    self.log.debug(f"Closed excess page. Total pages: {total_pages - 1}")
            else:
                self.available_pages.append(page)


class TurnstileAPIServer:
    HTML_TEMPLATE = """
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>Turnstile Solver</title>
        <script
          src="https://challenges.cloudflare.com/turnstile/v0/api.js?onload=onloadTurnstileCallback"
          async=""
          defer=""
        ></script>
      </head>
      <body>
        <!-- cf turnstile -->
      </body>
    </html>
    """

    def __init__(self, debug: bool = False):
        self.app = Quart(__name__)
        self.log = Logger()
        self.debug = debug
        self.page_pool = None
        self.browser = None
        self.context = None
        self.browser_args = [
            "--disable-blink-features=AutomationControlled",
        ]
        self._setup_routes()

    def _setup_routes(self) -> None:
        """Set up the application routes."""
        self.app.before_serving(self._startup)
        self.app.route('/turnstile', methods=['GET'])(self.process_turnstile)
        self.app.route('/')(self.index)

    async def _startup(self) -> None:
        """Initialize the browser and page pool on startup."""
        self.log.debug("Starting browser initialization...")
        try:
            await self._initialize_browser()
            self.log.success("Browser and page pool initialized successfully")
        except Exception as e:
            self.log.failure(f"Failed to initialize browser: {str(e)}")
            raise

    async def _initialize_browser(self) -> None:
        """Initialize the browser and create the page pool."""
        if self.debug:
            self.log.debug("Initializing browser with automation-resistant arguments...")

        playwright = await async_playwright().start()
        self.browser = await playwright.chromium.launch(
            headless=False,
            args=self.browser_args
        )
        self.context = await self.browser.new_context()
        self.page_pool = PagePool(
            self.context,
            debug=self.debug,
            log=self.log
        )
        await self.page_pool.initialize()

        if self.debug:
            self.log.debug(f"Browser and page pool initialized (min: {self.page_pool.min_size}, max: {self.page_pool.max_size})")

    async def _solve_turnstile(self, url: str, sitekey: str, invisible: bool = True) -> TurnstileAPIResult:
        """Solve the Turnstile challenge."""
        start_time = time.time()
        loader = Loader(desc="Solving captcha...", timeout=0.05)
        loader.start()

        page = await self.page_pool.get_page()
        try:
            if self.debug:
                self.log.debug(f"Starting Turnstile solve for URL: {url}")
                self.log.debug("Setting up page data and route")

            url_with_slash = url + "/" if not url.endswith("/") else url
            turnstile_div = f'<div class="cf-turnstile" data-sitekey="{sitekey}" data-theme="light"></div>'
            page_data = self.HTML_TEMPLATE.replace("<!-- cf turnstile -->", turnstile_div)

            await page.route(url_with_slash, lambda route: route.fulfill(body=page_data, status=200))
            await page.goto(url_with_slash)

            if self.debug:
                self.log.debug("Starting Turnstile response retrieval loop")

            max_attempts = 10
            attempts = 0
            while attempts < max_attempts:
                turnstile_check = await page.eval_on_selector(
                    "[name=cf-turnstile-response]", 
                    "el => el.value"
                )
                if turnstile_check == "":
                    if self.debug:
                        self.log.debug(f"Attempt {attempts + 1}: No Turnstile response yet")
                    
                    await page.evaluate("document.querySelector('.cf-turnstile').style.width = '70px'")
                    await page.click(".cf-turnstile")
                    await asyncio.sleep(0.5)
                    attempts += 1
                else:
                    turnstile_element = await page.query_selector("[name=cf-turnstile-response]")
                    if turnstile_element:
                        value = await turnstile_element.get_attribute("value")
                        elapsed_time = round(time.time() - start_time, 3)
                        
                        if self.debug:
                            self.log.debug(f"Turnstile response received: {value[:45]}...")
                        
                        self.log.message(
                            "Cloudflare",
                            f"Successfully solved captcha: {value[:45]}...",
                            start=start_time,
                            end=time.time()
                        )
                        
                        return TurnstileAPIResult(
                            result=value,
                            elapsed_time_seconds=elapsed_time
                        )
                    break

            self.log.failure("Failed to retrieve Turnstile value")
            return TurnstileAPIResult(
                result=None,
                status="failure",
                error="Max attempts reached without solution"
            )

        except Exception as e:
            self.log.failure(f"Error solving Turnstile: {str(e)}")
            return TurnstileAPIResult(
                result=None,
                status="error",
                error=str(e)
            )
        finally:
            loader.stop()
            if self.debug:
                self.log.debug("Clearing page state")
            await page.goto("about:blank")
            await self.page_pool.return_page(page)

    async def process_turnstile(self):
        """Handle the /turnstile endpoint requests."""
        url = request.args.get('url')
        sitekey = request.args.get('sitekey')
        invisible = request.args.get('invisible', 'false').lower() == 'true'

        if not url or not sitekey:
            self.log.warning("Missing required parameters: 'url' or 'sitekey'")
            return jsonify({
                "status": "error",
                "error": "Both 'url' and 'sitekey' are required"
            }), 400

        if self.debug:
            self.log.debug(f"Processing request for URL: {url}")
        try:
            result = await self._solve_turnstile(url=url, sitekey=sitekey, invisible=invisible)
            if self.debug:
                self.log.debug(f"Request completed with status: {result.status}")
            return jsonify(result.__dict__), 200 if result.status == "success" else 500
        except Exception as e:
            self.log.failure(f"Unexpected error processing request: {str(e)}")
            return jsonify({
                "status": "error",
                "error": str(e)
            }), 500

    async def index(self):
        """Serve the API documentation page."""
        return """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Turnstile Solver API</title>
            <script src="https://cdn.tailwindcss.com"></script>
        </head>
        <body class="bg-gray-100 min-h-screen flex items-center justify-center">
            <div class="bg-white p-8 rounded-lg shadow-md max-w-2xl w-full">
                <h1 class="text-3xl font-bold mb-6 text-center text-blue-600">Welcome to Turnstile Solver API</h1>

                <p class="mb-4 text-gray-700">To use the turnstile service, send a GET request to 
                   <code class="bg-gray-200 px-2 py-1 rounded">/turnstile</code> with the following query parameters:</p>

                <ul class="list-disc pl-6 mb-6 text-gray-700">
                    <li><strong>url</strong>: The URL where Turnstile is to be validated</li>
                    <li><strong>sitekey</strong>: The site key for Turnstile</li>
                    <li><strong>invisible</strong>: Optional. Set to true if the Turnstile is invisible.</li>
                </ul>

                <div class="bg-gray-200 p-4 rounded-lg mb-6">
                    <p class="font-semibold mb-2">Example usage:</p>
                    <code class="text-sm break-all">/turnstile?url=https://example.com&sitekey=sitekey&invisible=true</code>
                </div>

                <div class="bg-blue-100 border-l-4 border-blue-500 p-4 mb-6">
                    <p class="text-blue-700">This project is inspired by 
                       <a href="https://github.com/Body-Alhoha/turnaround" class="text-blue-600 hover:underline">Turnaround</a> 
                       and is currently maintained by 
                       <a href="https://github.com/Theyka" class="text-blue-600 hover:underline">Theyka</a> 
                       and <a href="https://github.com/sexfrance" class="text-blue-600 hover:underline">Sexfrance</a>.</p>
                </div>
            </div>
        </body>
        </html>
        """

    def create_app(self):
        """Create and configure the application instance."""
        return self.app

if __name__ == "__main__":
    server = TurnstileAPIServer(debug=True)
    app = server.create_app()
    app.run()

# Credits for the changes: github.com/sexfrance
# Credit for the original script: github.com/Theyka
