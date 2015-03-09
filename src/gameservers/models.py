
from datetime import timedelta
from ast import literal_eval

from django.db import models
from django.utils import timezone


DAY_NAMES = [
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday',
    'Sunday',
]


class Server(models.Model):
    title = models.CharField(max_length=255)
    game_url = models.CharField(max_length=255)
    site_url = models.URLField(blank=True)

    last_updated = models.DateTimeField(auto_now=True, default=timezone.now)
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
        ordering = ['-players_current', 'title']

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
            int(stats['players__avg'] or 0),
            stats['players__min'] or 0,
            stats['players__max'] or 0,
        )

    def measure_weekdays(self, days=7):
        weekdays = []
        history = self.get_stats_history(days=days)
        for i, day in enumerate(DAY_NAMES):
            # HACK: Fuck django for defaulting the start of a week to sunday...
            # Oh, and it's not even zero indexed like the rest of the python
            # community! So much fun with number juggling...
            i += 2
            if i > 7: i = 1
            # NOTE: using __week_day is dependant on pytz
            tmp = history.filter(created__week_day=i)
            avg = tmp.aggregate(models.Avg('players'))['players__avg'] or 0
            weekdays.append((day, int(avg)))

        return weekdays

    def update_stats(self, player_count=0):
        self.players_current = player_count

        tmp = self.measure_players(days=31)
        self.players_avg, self.players_min, self.players_max = tmp

        tmp = self.measure_weekdays()
        self.averages_for_weekdays = ','.join([str(i) for day, i in tmp])

    def get_averages_for_weekdays(self):
        tmp = literal_eval(self.averages_for_weekdays)
        return zip(DAY_NAMES, tmp)

class ServerHistory(models.Model):
    server = models.ForeignKey(Server)
    created = models.DateTimeField(default=timezone.now)
    players = models.PositiveIntegerField(default=0)

    class Meta:
        ordering = ['-created', 'server']

    def __str__(self):
        return 'History for {} at {}.'.format(self.server, self.created)

