
from multiprocessing import Pool
import socket
import struct

import re
import requests
from bs4 import BeautifulSoup


def poll_ss13_server(host, port, timeout=30):
    # Thanks to /u/headswe for showing how to poll servers.
    # Source: http://www.reddit.com/r/SS13/comments/31b5im/a_bunch_of_graphs_for_all_servers/cq11nld
    print 'polling:', host, port
    query = '\x00\x83{0}\x00\x00\x00\x00\x00?players\x00'.format(
        struct.pack('>H', 14) # 14 = 6 null bytes + len(?players)
    )

    try:
        sock = socket.create_connection((host, port), timeout=timeout)
    except socket.timeout:
        return

    try:
        sock.sendall(query)
        response = sock.recv(1024)
    except socket.timeout:
        response = ''

    sock.close()
    print 'done:', host, port
    if len(response) < 1:
        return
    else:
        if not response[:5] == '\x00\x83\x00\x05\x2a':
            return
        tmp = struct.unpack('f', response[5:9])
        return (int(tmp[0]), host, port)


class ServerPoller(object):
    def __init__(self, timeout=10):
        self.timeout = timeout
        self.workers = 5

    def run(self):
        targets = self._get_servers()
        servers = self._poll_servers(targets)
        return servers

    def _get_servers(self):
        return [
            ('baystation12.net', 8000),
            ('8.8.4.4', 3333),
            ('ss13.lljk.net', 26100),
            ('204.152.219.158', 3333),
        ]

    def _poll_servers(self, targets):
        pool = Pool(processes=self.workers)
        results = []
        for (host, port) in targets:
            future = pool.apply_async(poll_ss13_server, (host, port, timeout))
            results.append(future)

        pool.close()
        pool.join()

        servers = []
        for future in results:
            server = self._handle_future(future)
            if server:
                servers.append(server)
        return servers

    def _handle_future(self, future):
        tmp = future.get()
        if not tmp:
            return
        (players, host, port) = tmp
        server = dict(
            title = host,
            game_url = 'byond://{}:{}'.format(host, port),
            site_url = '',
            player_count = players,
        )
        return server


class ServerScraper(object):
    def __init__(self):
        self.url = 'http://www.byond.com/games/exadv1/spacestation13'
        # TODO: Better regexp that can't be spoofed by server names
        self.PLAYERS = re.compile('Logged in: (\d+) player')

    def run(self):
        '''Run the parser and return a neat list of dicts containing server data.'''
        raw_data = self._download_data()
        servers = self._parse_data(raw_data)
        return servers

    def _download_data(self):
        '''Download raw data, either from local file or a web page.'''
        if self.url.startswith('http://') or self.url.startswith('https://'):
            raw_data = requests.get(self.url).text.strip()
        else:
            # HACK: In case of local testing or debugging, since requests can't
            # handle local files.
            with open(self.url, 'r') as f:
                raw_data = f.read().strip()
        return raw_data

    def _parse_data(self, raw_data):
        '''Parse the raw data and run through all servers.'''
        servers = []
        soup_data = BeautifulSoup(raw_data)
        for server_data in soup_data.find_all('div', 'live_game_status'):
            server = self._parse_server_data(server_data)
            if server:
                servers.append(server)

        return servers

    def _parse_server_data(self, data):
        '''Parse the individual parts of each server.'''
        try:
            title = data.find('b').get_text().splitlines()[0].strip().encode('utf-8')
        except AttributeError:
            # HACK: I think this happends because the raw data was incomplete.
            # No complete data, no server update.
            return None
        game_url = data.find('span', 'smaller').text.encode('utf-8')

        tmp = data.find('a')
        site_url = None
        # Default means the server hasn't set a custom site url
        if tmp and not tmp.text == 'Default':
            try:
                site_url = tmp['href'].encode('utf-8')
                # Handle some funky servers...
                if site_url == 'http://':
                    site_url = ''
            except KeyError:
                # Sometimes there's a <a> tag without a href attribute
                pass

        tmp = data.text
        player_count = 0
        if tmp.find('No players.') == -1:
            data = self.PLAYERS.search(tmp)
            player_count = int(data.group(1))

        server = dict(
            title = title,
            game_url = game_url,
            site_url = site_url,
            player_count = player_count,
        )
        return server

