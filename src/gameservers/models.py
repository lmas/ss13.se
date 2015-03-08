
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


class ServerHistory(models.Model):
    server = models.ForeignKey(Server)
    created = models.DateTimeField(default=timezone.now)
    players = models.PositiveIntegerField(default=0)

    class Meta:
        ordering = ['-created', 'server']

    def __str__(self):
        return 'History for {} at {}.'.format(self.server, self.created)

