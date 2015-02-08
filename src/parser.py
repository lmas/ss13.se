#!/usr/bin/env python

import re
import logging

import requests
from bs4 import BeautifulSoup


URL = 'http://www.byond.com/games/exadv1/spacestation13'
PLAYER_COUNT = re.compile('Logged in: (\d+) player')


logging.basicConfig(
    format = '%(asctime)s %(levelname)s %(message)s',
    level = logging.INFO,
)


class ServerParser(object):
    def __init__(self):
        self.url = URL

    def run(self):
        '''Run the parser and return a neat list of dicts containing server data.'''
        raw_data = self._download_data()
        servers = self._parse_data(raw_data)
        return servers

    def _download_data(self):
        '''Download raw data, either from local file or a web page.'''
        logging.info('Downloading data from {} ...'.format(self.url))
        if self.url.startswith('http://') or self.url.startswith('https://'):
            raw_data = requests.get(self.url).text.strip()
        else:
            # HACK: In case of local testing or debugging, since requests can't
            # handle local files.
            logging.debug('Opening local file...')
            with open(self.url, 'r') as f:
                raw_data = f.read().strip()
        return raw_data

    def _parse_data(self, raw_data):
        '''Parse the raw data and run through all servers.'''
        logging.info('Parsing raw data...')
        servers = []
        soup_data = BeautifulSoup(raw_data)
        for server_data in soup_data.find_all('div', 'live_game_status'):
            server = self._parse_server_data(server_data)
            servers.append(server)

        logging.info('Number of servers parsed: {}'.format(len(servers)))
        return servers

    def _parse_server_data(self, data):
        '''Parse the individual parts of each server.'''
        title = data.find('b').text.strip()
        game_url = data.find('span', 'smaller').text

        tmp = data.find('a')
        site_url = None
        # Default means the server hasn't set a custom site url
        if tmp and not tmp.text == 'Default':
            try:
                site_url = tmp['href']
                # Handle some funky servers...
                if site_url == 'http://':
                    site_url = ''
            except KeyError:
                # Sometimes there's a <a> tag without a href attribute
                pass

        tmp = data.text
        player_count = 0
        if tmp.find('No players.') == -1:
            data = PLAYER_COUNT.search(tmp)
            player_count = int(data.group(1))

        server = dict(
            title = title,
            game_url = game_url,
            site_url = site_url,
            player_count = player_count,
        )
        return server


if __name__ == '__main__':
    parser = ServerParser()
    parser.url = './dump.html' # Use a local file instead when testing
    servers = parser.run()
    for tmp in servers:
        print '{}\nPlayers: {}\n'.format(tmp['title'], tmp['player_count'])

