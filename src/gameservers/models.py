
import datetime

from django.db import models
from django.utils import timezone

import redis


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


class PlayerHistory(object):
    def __init__(self, redis_settings=dict(host='localhost', port=6379, db=0)):
        self.redis = redis.StrictRedis(**redis_settings)

        # 2688 = 4 times per hour * 24 hours * 7 days * 4 weeks
        self.max_points = 2688

    def add_point(self, server, time, players):
        '''Add a new point in the player history.'''
        self.redis.lpush(server, '{},{}'.format(time, players))

    def trim_points(self, server):
        '''Trim away too old points and servers in the player history.'''
        self.redis.ltrim(server, 0, self.max_points)
        # let the list expire after a week without updates
        self.redis.expire(server, 604800)

    def get_points(self, server):
        '''Get a range of points from the player history.'''
        # TODO: not memory efficient, turn into a generator instead?
        points = []
        for tmp in self.redis.lrange(server, 0, self.max_points):
            time, players = tmp.split(',')
            time, players = datetime.datetime.fromtimestamp(float(time)), int(players)
            points.append((time, players))
        return points

