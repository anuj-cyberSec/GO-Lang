from __future__ import annotations

import concurrent.futures
import sys
import time
import urllib.error
import urllib.request


def check_url(url: str, timeout_s: float = 3.0) -> tuple[str, int | None, float, str | None]:
    start = time.perf_counter()
    try:
        with urllib.request.urlopen(url, timeout=timeout_s) as response:
            status = response.getcode()
            return url, status, time.perf_counter() - start, None
    except (urllib.error.URLError, TimeoutError) as exc:
        return url, None, time.perf_counter() - start, str(exc)


def main() -> None:
    urls = sys.argv[1:]
    if not urls:
        urls = [
            "https://example.com",
            "https://httpbin.org/status/200",
            "https://httpbin.org/status/503",
            "https://httpbin.org/delay/1",
            "https://httpbin.org/delay/2",
        ]

    max_workers = 4
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(check_url, url) for url in urls]
        for future in concurrent.futures.as_completed(futures):
            url, status, duration, err = future.result()
            if err:
                print(f"{url:<35} error={err} duration={duration:.3f}s")
            else:
                print(f"{url:<35} status={status} duration={duration:.3f}s")


if __name__ == "__main__":
    main()
