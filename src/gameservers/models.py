
from datetime import timedelta
import calendar

from django.db import models
from django.utils import timezone


class Server(models.Model):
    title = models.CharField(max_length=255)
    game_url = models.CharField(max_length=255)
    site_url = models.URLField(blank=True)
    current_players = models.PositiveIntegerField(default=0, editable=False)
    last_updated = models.DateTimeField(auto_now=True, default=timezone.now)

    class Meta:
        ordering = ['-current_players', 'title']

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

    def get_history_stats(self, days=7):
        return ServerHistory.objects.filter(
            server=self,
            created__gte=timezone.now() - timedelta(days=days)
        )

    def calc_player_stats(self, days=7):
        history = self.get_history_stats(days=days)
        return history.aggregate(
            models.Avg('players'),
            models.Min('players'),
            models.Max('players'),
        )

    def weekday_averages(self):
        weekdays = []
        history = self.get_history_stats(days=7)
        for i, day in enumerate(calendar.day_name):
            # HACK: do some number juggling to convert from calendar to django,
            # because SOMEONE didn't bother to follow THE FUCKING STANDARD
            #
            # calendar is zero indexed, first day of week defaults to monday
            # (monday = 0, tuesday = 1 etc.)
            #
            # django isn't zero indexed, first day of week defaults to sunday
            # (sunday = 1, monday = 2 etc.)
            i += 2
            if i > 7: i = 1
            # NOTE: using __week_day is dependant on pytz
            tmp = history.filter(created__week_day=i)
            avg = tmp.aggregate(models.Avg('players'))['players__avg'] or 0
            weekdays.append((day, int(avg)))
        return weekdays


class ServerHistory(models.Model):
    server = models.ForeignKey(Server)
    created = models.DateTimeField(default=timezone.now)
    players = models.PositiveIntegerField(default=0)

    class Meta:
        ordering = ['-created', 'server']

    def __str__(self):
        return 'History for {} at {}.'.format(self.server, self.created)

