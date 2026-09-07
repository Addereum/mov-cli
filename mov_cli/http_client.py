from __future__ import annotations
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from typing import Literal, Dict, Optional

    from httpx import Response

import httpx
import warnings
from http.cookiejar import MozillaCookieJar
from devgoldyutils import LoggerAdapter, Colours

from .utils import hide_ip
from .utils.cookies import get_cookie_file
from .logger import mov_cli_logger
from .errors import SiteMaybeBlockedError
from .config import Config

__all__ = ("HTTPClient",)

class HTTPClient():
    def __init__(
        self, 
        headers: Optional[Dict[str, str]] = None, 
        timeout: int = 15, 
        hide_ip: bool = True,
        proxy: Optional[Dict[str, str]] = None,
        config: Optional[Config] = None
    ) -> None:
        self.hide_ip = hide_ip
        self.headers = headers or {}
        self.proxy = proxy

        self.logger = LoggerAdapter(mov_cli_logger, prefix = self.__class__.__name__)

        if self.proxy:
            self.logger.debug(f"Using proxy -> {self.proxy}")

        cookies = None
        
        if config is None:
            config = Config()
            
        cookie_file = get_cookie_file(config)
        if cookie_file is not None:
            try:
                cj = MozillaCookieJar(cookie_file)
                cj.load(ignore_discard=True, ignore_expires=True)
                cookies = cj
            except Exception as e:
                self.logger.warning(f"Failed to load cookies from {cookie_file}: {e}")

        self.__httpx_client = httpx.Client(
            timeout = timeout, 
            cookies = cookies,
            proxy = self.proxy
        )

        super().__init__()

    def close(self) -> None:
        """Close the underlying httpx client and release connections."""
        self.__httpx_client.close()

    def __enter__(self):
        return self

    def __exit__(self, *args):
        self.close()

    def request(
        self, 
        method: Literal["GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"],
        url: str, 
        params: Optional[Dict[str, str]] = None, 
        headers: Optional[Dict[str, str]] = None, 
        include_default_headers: bool = False, 
        redirect: bool = False, 
        **kwargs
    ) -> Response:
        """Performs a request with httpx and returns `httpx.Response`."""
        if headers is None:
            headers = {}

        if include_default_headers is True:

            if headers.get("Referer") is None:
                headers.update({"Referer": url})

            headers.update(self.headers)

        try:
            self.logger.debug(
                Colours.ORANGE.apply(method.upper()) + f" -> {hide_ip(url, self.hide_ip)}"
            )

            response = self.__httpx_client.request(
                method = method, 
                url = url, 
                params = params, 
                headers = headers, 
                follow_redirects = redirect, 
                **kwargs
            )

            if response.is_error:
                self.logger.debug(
                    f"{method.upper()} request to '{response.url}' {Colours.RED.apply('failed!')} ({response})"
                )

            return response

        except httpx.ConnectError as e:
            error_str = str(e)

            if any(ssl_hint in error_str for ssl_hint in (
                "[SSL: CERTIFICATE_VERIFY_FAILED]",
                "[SSL]",
                "CERTIFICATE_VERIFY_FAILED",
                "SSLError",
                "SSLCertVerificationError",
            )):
                raise SiteMaybeBlockedError(url, e)

            raise e

        except httpx.ConnectTimeout as e:
            self.logger.warning(
                f"Connection to '{hide_ip(url, self.hide_ip)}' timed out. "
                "The site may be down or your connection may be slow."
            )
            raise e

    def get(
        self, 
        url: str, 
        headers: Dict[str, str] = {}, 
        include_default_headers: bool = True, 
        redirect: bool = False, 
        **kwargs
    ) -> Response:
        """Performs a GET request and returns httpx.Response.

        .. deprecated:: 4.4
            Use :meth:`HTTPClient.request` instead.
        """
        warnings.warn(
            "HTTPClient.get() is deprecated since v4.4. Use HTTPClient.request() instead.",
            DeprecationWarning,
            stacklevel = 2
        )
        return self.request(
            "GET", 
            url = url, 
            headers = headers, 
            include_default_headers = include_default_headers, 
            redirect = redirect,
            **kwargs
        )

    def post(
        self, 
        url: str,
        data: dict = {},
        json: dict = {}, 
        headers: dict = {}, 
        include_default_headers: bool = True, 
        redirect: bool = False, 
        **kwargs
    ) -> Response:
        """Performs a POST request and returns httpx.Response.

        .. deprecated:: 4.4
            Use :meth:`HTTPClient.request` instead.
        """
        warnings.warn(
            "HTTPClient.post() is deprecated since v4.4. Use HTTPClient.request() instead.",
            DeprecationWarning,
            stacklevel = 2
        )
        return self.request(
            "POST", 
            url = url, 
            data = data, 
            json = json, 
            headers = headers, 
            include_default_headers = include_default_headers, 
            redirect = redirect,
            **kwargs
        )

    def set_cookies(self, cookies: dict) -> None:
        """Sets cookies."""
        self.__httpx_client.cookies = cookies