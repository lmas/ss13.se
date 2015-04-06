
from multiprocessing import Pool
import socket
import struct

import re
import requests
from bs4 import BeautifulSoup

from .models import PrivateServer


def poll_ss13_server(server, timeout=30):
    '''Try connect to a SS13 server and get a player count from it.'''
    # Thanks to /u/headswe for showing how to poll servers.
    # Source: http://www.reddit.com/r/SS13/comments/31b5im/a_bunch_of_graphs_for_all_servers/cq11nld
    addr = (server.host, server.port)
    query = '\x00\x83{0}\x00\x00\x00\x00\x00?players\x00'.format(
        struct.pack('>H', 14) # 14 = 6 null bytes + len(?players)
    )

    try:
        sock = socket.create_connection(addr, timeout=timeout)
        sock.sendall(query)
        response = sock.recv(1024)
        sock.close()
        assert(len(response) >= 9)
        assert(response[:5] == '\x00\x83\x00\x05\x2a')
        players = int(struct.unpack('f', response[5:9])[0])
        assert(players >= 0)
        return players, server
    except (socket.timeout, AssertionError) as e:
        try:
            sock.close()
        except UnboundLocalError:
            pass
        return -1, server


class ServerPoller(object):
    '''Manually poll hidden/private servers for their stats.'''
    def __init__(self, timeout=10):
        self.timeout = timeout
        self.workers = 5

    def run(self):
        '''Run the poller and return a nice list of dicts containing server data.'''
        targets = self._get_servers()
        servers = self._poll_servers(targets)
        return servers

    def _get_servers(self):
        '''Grab all private servers that's been manually activated.'''
        servers = PrivateServer.objects.filter(active=True)
        return servers

    def _poll_servers(self, targets):
        '''Poll each server in the targets list and try get it's stats.'''
        pool = Pool(processes=self.workers)
        results = []
        for server in targets:
            future = pool.apply_async(poll_ss13_server, (server, self.timeout))
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
        '''Check if a poll of a server succeeded and return a data dict for it.'''
        players, server = future.get()
        if players == -1:
            # Couldn't get a proper update from the server when polling it,
            # consider it offline for now
            return
        server = dict(
            title = server.title,
            game_url = 'byond://{}:{}'.format(server.host, server.port),
            site_url = server.site_url,
            player_count = players,
        )
        return server


class ServerScraper(object):
    '''Scrape the Byond server page for server stats.'''

    def __init__(self):
        self.url = 'http://www.byond.com/games/exadv1/spacestation13'
        # TODO: Better regexp that can't be spoofed by server names
        self.PLAYERS = re.compile('Logged in: (\d+) player')

    def run(self):
        '''Run the scraper and return a neat list of dicts containing server data.'''
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

