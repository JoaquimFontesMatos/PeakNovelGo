import cloudscraper
import requests
import random

PROXY_SCRAPE_URL = "https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&proxy_format=protocolipport&format=text"

USER_AGENTS = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:48.0) Gecko/20100101 Firefox/48.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
    "Mozilla/5.0 (Windows NT 6.3; Trident/7.0; AS3366; AS3366; rv:11.0) like Gecko",
]

class Request:
    def __init__(self):
        self.s = cloudscraper.create_scraper()
        self.s.headers.update({
            "User-Agent": random.choice(USER_AGENTS),  # Randomize User-Agent
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.5",
            "Referer": "https://novel-bin.net/",
            "Connection": "keep-alive",  # Add Connection header
            "Upgrade-Insecure-Requests": "1",  # Add Upgrade-Insecure-Requests header
        })
        self.proxies = self.get_proxies()

    def get_proxies(self):
        """Fetch proxies from ProxyScrape and return a list."""
        try:
            response = requests.get(PROXY_SCRAPE_URL)
            if response.status_code == 200:
                proxy_list = response.text.strip().split("\n")
                return [p.strip() for p in proxy_list if p.strip()]
            else:
                #print("❌ Failed to fetch proxies!")
                return []
        except Exception as e:
            #print(f"❌ Error fetching proxies: {e}")
            return []

    def get_random_proxy(self):
        """Select a random proxy and format it properly."""
        if not self.proxies:
            #print("⚠️ No proxies available! Using direct connection.")
            return None

        raw_proxy = random.choice(self.proxies)
        #print(f"🌍 Using proxy: {raw_proxy}")

        try:
            protocol, ip_port = raw_proxy.split("://")
            return {protocol: ip_port}
        except ValueError:
            #print(f"⚠️ Invalid proxy format: {raw_proxy}")
            return None

    def request_with_proxy(self, method, url, data=None):
        """Try request with a proxy, retry if blocked."""
        for _ in range(5):  # Retry up to 5 times with different proxies
            proxy = self.get_random_proxy()
            if proxy is None:
                continue  # Skip if proxy is invalid

            try:
                self.s.headers.update({"User-Agent": random.choice(USER_AGENTS)})  # Randomize User-Agent per request
                response = self.s.request(method, url, data=data, proxies=proxy, timeout=10)

                # If Cloudflare blocks, retry
                if response.status_code in [403, 503, 429]:
                    #print(f"🚫 Blocked! Retrying with a new proxy... ({response.status_code})")
                    continue  # Try another proxy

                return response  # Valid response found
            except requests.exceptions.RequestException as e:
                continue

        #print("❌ All proxies failed! Returning None.")
        return None  # No working proxies

    def get(self, url):
        return self.request_with_proxy("GET", url)

    def post(self, url, data):
        return self.request_with_proxy("POST", url, data)
