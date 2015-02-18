from django.db import models

class Server(models.Model):
    title = models.CharField(max_length=255)
    game_url = models.URLField()
    site_url = models.URLField(blank=True)

    class Meta:
        ordering = ['title']

    def __str__(self):
        return self.title

class Population(models.Model):
    timestamp = models.DateTimeField(auto_now_add=True)
    server = models.ForeignKey(Server)
    players = models.PositiveIntegerField()

    class Meta:
        ordering = ['-timestamp', 'server']

    def __str__(self):
        return '{} {}'.format(self.timestamp, self.server.title)

