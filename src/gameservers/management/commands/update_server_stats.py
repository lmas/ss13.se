#!/usr/bin/env python

import re

from django.core.management.base import BaseCommand
from django.utils import timezone

from gameservers.models import Server, ServerHistory

import requests
from bs4 import BeautifulSoup


URL = 'http://www.byond.com/games/exadv1/spacestation13'
# TODO: Better regexp that can't be spoofed by server names
PLAYER_COUNT = re.compile('Logged in: (\d+) player')


class ServerScraper(object):
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
            data = PLAYER_COUNT.search(tmp)
            player_count = int(data.group(1))

        server = dict(
            title = title,
            game_url = game_url,
            site_url = site_url,
            player_count = player_count,
        )
        return server


from multiprocessing import Pool
import socket
import struct


def poll_ss13_server(host, port, timeout=10):
    # Thanks to /u/headswe for showing how to poll servers.
    # Source: http://www.reddit.com/r/SS13/comments/31b5im/a_bunch_of_graphs_for_all_servers/cq11nld
    print 'polling:', host, port
    cmd = '?players'
    query = '\x00\x83{0}\x00\x00\x00\x00\x00{1}\x00'.format(
        struct.pack('>H', len(cmd) + 6), cmd
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
    def __init__(self, timeout=30):
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
            future = pool.apply_async(poll_ss13_server, (host, port))
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


class Command(BaseCommand):
    help = 'Update history stats for all ss13 servers.'

    def _update_stats(self, server, players, time):
        # Create a new record in the history
        history = ServerHistory(server=server, players=players)

        # Update "live stats"
        server.update_stats(players, time=time)
        server.save()

        return history

    def handle(self, *args, **kwargs):
        servers = []

        # Poll private servers not on the Byond server page
        # Prioritize these servers and make them be handled first.
        poller = ServerPoller()
        servers.extend(poller.run())

        # Now we scrape servers of the Byond server page
        # We have less control of these, so they become less prioritized
        # (anyone could make a clone of a server with zero players on,
        # just to fuck with the cloned server's stats and graphs).
        parser = ServerScraper()
        #parser.url = './dump.html' # Use a local file instead when testing
        servers.extend(parser.run())

        servers_handled = []
        new_items = []
        now = timezone.now()

        for data in servers:
            # Prevent empty servers with identical names to other, active servers
            # from fucking with the history
            if data['title'] in servers_handled:
                continue
            else:
                servers_handled.append(data['title'])

            # Grab the correct server, making sure it has updated links
            server, created = Server.objects.update_or_create(
                title=data['title'],
                defaults= dict(
                    game_url=data['game_url'],
                    site_url=data['site_url'] or '',
                )
            )

            tmp = self._update_stats(server, data['player_count'], now)
            new_items.append(tmp)

        # Make sure to update servers not available on the page.
        for server in Server.objects.exclude(last_updated__exact=now):
            tmp = self._update_stats(server, 0, None)
            new_items.append(tmp)

        ServerHistory.objects.bulk_create(new_items)
        Server.remove_old_servers()

