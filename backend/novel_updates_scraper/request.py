import random

from playwright.sync_api import sync_playwright, TimeoutError as PlaywrightTimeoutError

USER_AGENTS = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
]


class PlaywrightScraper:
    def __init__(self):
        self.user_agent = random.choice(USER_AGENTS)

    def scrape(self, url):
        with sync_playwright() as p:
            # Launch the browser (use headless=False for debugging)
            browser = p.chromium.launch(headless=True)
            context = browser.new_context(user_agent=self.user_agent)

            # Open a new page
            page = context.new_page()

            try:
                # Navigate to the URL
                response = page.goto(url, wait_until="domcontentloaded", timeout=10000)

                if not response:
                    return None

                page.wait_for_selector("body", state="attached")

                # Check if the response is successful
                if not response or response.status != 200:
                    if response and response.status == 404:
                        return "Page not found"

                    return None

                # Return the page content
                return page.content()
            except PlaywrightTimeoutError:
                return None
            finally:
                browser.close()
