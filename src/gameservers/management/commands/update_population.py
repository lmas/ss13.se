#!/usr/bin/env python

import re

from django.core.management.base import BaseCommand

from gameservers.models import Server, ServerHistory

import requests
from bs4 import BeautifulSoup


URL = 'http://www.byond.com/games/exadv1/spacestation13'
# TODO: Better regexp that can't be spoofed by server names
PLAYER_COUNT = re.compile('Logged in: (\d+) player')


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
            title = data.find('b').get_text().splitlines()[0].strip()
        except AttributeError:
            # HACK: I think this happends because the raw data was incomplete.
            # No complete data, no server update.
            return None
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


class Command(BaseCommand):
    help = 'Update population stats for all ss13 servers.'

    def handle(self, *args, **kwargs):
        parser = ServerParser()
        #parser.url = './dump.html' # Use a local file instead when testing
        servers = parser.run()
        servers_handled = []
        new_items = []

        for data in servers:
            # Prevent empty servers with identical names to other, active servers
            # from fucking with the history
            if data['title'] in servers_handled:
                continue
            else:
                servers_handled.append(data['title'])

            server, created = Server.objects.update_or_create(
                title=data['title'],
                defaults= dict(
                    game_url=data['game_url'],
                    site_url=data['site_url'] or '',
                    current_players=data['player_count'],
                )
            )

            history = ServerHistory(server=server, players=data['player_count'])
            new_items.append(history)

        ServerHistory.objects.bulk_create(new_items)
        Server.remove_old_servers()

