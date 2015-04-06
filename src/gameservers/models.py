
from __future__ import unicode_literals
from datetime import timedelta
from ast import literal_eval

from django.db import models
from django.utils import timezone
from django.utils.encoding import python_2_unicode_compatible


DAY_NAMES = [
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday',
    'Sunday',
]


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
    averages_for_weekdays = models.CommaSeparatedIntegerField(
        max_length=50,
        editable=False,
        default='',
    )

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

    def measure_weekdays(self, days=7):
        history = self.get_stats_history(days=days)
        weekdays = []
        # Why can't it be zero indexed like the rest of the fucking community...
        for day in range(1, 8):
            tmp = history.filter(created__week_day=day)
            avg = tmp.aggregate(models.Avg('players'))['players__avg'] or 0
            weekdays.append(int(round(avg, 0)))
        # HACK: Since django's __week_day starts on a sunday (amurican suckers)
        # we have to move sunday (at the start) to the end of the list
        weekdays.insert(len(weekdays), weekdays.pop(0))
        return weekdays

    def update_stats(self, player_count=0, time=None):
        # TODO: default to setting current time
        if time:
            self.last_updated = time

        self.players_current = player_count

        tmp = self.measure_players(days=31)
        self.players_avg, self.players_min, self.players_max = tmp

        tmp = self.measure_weekdays(days=31)
        self.averages_for_weekdays = ','.join([str(i) for i in tmp])

    def get_averages_for_weekdays(self):
        try:
            tmp = literal_eval(self.averages_for_weekdays)
        except SyntaxError:
            tmp = [0,0,0,0,0,0,0]
        return zip(DAY_NAMES, tmp)

@python_2_unicode_compatible
class ServerHistory(models.Model):
    server = models.ForeignKey(Server)
    created = models.DateTimeField(default=timezone.now)
    players = models.PositiveIntegerField(default=0)

    class Meta:
        ordering = ['-created', 'server']

    def __str__(self):
        return 'History for {} at {}.'.format(self.server.title, self.created)

