
from __future__ import unicode_literals
from datetime import timedelta
from ast import literal_eval

from django.db import models
from django.utils import timezone
from django.utils.encoding import python_2_unicode_compatible


@python_2_unicode_compatible
class PrivateServer(models.Model):
    title = models.CharField(max_length=255)
    site_url = models.URLField(blank=True)
    host = models.CharField(max_length=255)
    port = models.PositiveIntegerField()
    active = models.BooleanField(default=False)

    class Meta:
        ordering = ['-active', 'title']

    def __str__(self):
        return self.title

    @staticmethod
    def deactivate_server(server):
        try:
            tmp = PrivateServer.objects.get(
                title=server.title,
                site_url=server.site_url,
            )
        except PrivateServer.DoesNotExist:
            return
        tmp.active = False
        tmp.save()


@python_2_unicode_compatible
class Server(models.Model):
    title = models.CharField(max_length=255)
    game_url = models.CharField(max_length=255)
    site_url = models.URLField(blank=True)

    last_updated = models.DateTimeField(default=timezone.now, editable=False)
    players_current = models.PositiveIntegerField(default=0, editable=False)
    players_avg = models.PositiveIntegerField(default=0, editable=False)
    players_min = models.PositiveIntegerField(default=0, editable=False)
    players_max = models.PositiveIntegerField(default=0, editable=False)

    class Meta:
        ordering = ['-players_current', '-last_updated', 'title']

    def __str__(self):
        return self.title

    @staticmethod
    def remove_old_servers():
        '''Remove servers that hasn't been updated for over a week.'''
        now = timezone.now()
        for server in Server.objects.all():
            delta = now - server.last_updated
            if delta.days >= 7:
                server.delete()
                PrivateServer.deactivate_server(server)

    def get_stats_history(self, days=7):
        return ServerHistory.objects.filter(
            server=self,
            created__gte=timezone.now() - timedelta(days=days)
        )

    def measure_players(self, days=7):
        history = self.get_stats_history(days=days)
        stats = history.aggregate(
            models.Avg('players'),
            models.Min('players'),
            models.Max('players'),
        )
        return (
            int(round(stats['players__avg'] or 0, 0)),
            stats['players__min'] or 0,
            stats['players__max'] or 0,
        )

    def update_stats(self, player_count=0, time=None):
        # TODO: default to setting current time
        if time:
            self.last_updated = time

        self.players_current = player_count

        tmp = self.measure_players(days=31)
        self.players_avg, self.players_min, self.players_max = tmp


@python_2_unicode_compatible
class ServerHistory(models.Model):
    server = models.ForeignKey(Server)
    created = models.DateTimeField(default=timezone.now)
    players = models.PositiveIntegerField(default=0)

    class Meta:
        ordering = ['-created', 'server']

    def __str__(self):
        return 'History for {} at {}.'.format(self.server.title, self.created)

