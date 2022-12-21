import logging
import requests
from random import choice
from lxml import etree

_global_agent = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
]


class HttpServer:
    def __init__(self, base: str = "", header: dict = {}):
        self._base = base
        self._header = header
        # 默认添加浏览器header
        self._header["user-agent"] = choice(_global_agent)
        self._proxy = {}

    def set_proxy(self, proxy: dict):
        self._proxy = proxy
        return self

    def set_header(self, key: str, value: str):
        self._header[key] = value
        return self

    def _get(self, url: str):
        _url = "{}{}".format(self._base, url)
        logging.info("request url {}".format(_url))
        response = requests.get(_url, headers=self._header, proxies=self._proxy)
        if response.status_code != 200:
            logging.error("request code is not 200 res {}".format(response.text))
            return None
        return response

    def get_etree(self, url: str) -> etree:
        _response = self._get(url)
        _res = ''
        if _response is not None:
            _res = _response.text
        return etree.HTML(_res, etree.HTMLParser())


if __name__ == '__main__':
    pass
    # server = HttpServer("https://dmhy.anoneko.com", {})
    # server.set_proxy({'http': 'http://192.168.1.1:7890', 'https': 'http://192.168.1.1:7890'})
    # res = server.get("/topics/view/624838_pop_pipi_Pop_Team_Epic_S2_12_END_1080P_AVC_AAC_MP4.html")
    # with open("tmp.txt", "w", encoding="utf-8") as f:
    #     f.write(res)
    # with open("tmp.txt", "r", encoding="utf-8") as f:
    #     data = f.read()
    #     html = etree.HTML(data, etree.HTMLParser())
    #     # 遍历所有的表格
    #     for resource in html.xpath('//div[@id="resource-tabs"]/div[@id="tabs-1"]/p'):
    #         title = resource.xpath('strong/text()')
    #         url = resource.xpath('a/@href')
    #         url_title = resource.xpath('a/text()')
    #         if len(title) > 0:
    #             print("{}-{}-{}".format(title[0], url[0], url_title[0]))
