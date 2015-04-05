#!/usr/bin/env python

from django.core.management.base import BaseCommand
from django.utils import timezone

from gameservers.models import Server, ServerHistory

from gameservers.data_sources import ServerPoller, ServerScraper

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

