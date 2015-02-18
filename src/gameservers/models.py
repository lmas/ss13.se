from django.db import models

class Server(models.Model):
    title = models.CharField(max_length=255)
    game_url = models.URLField()
    site_url = models.URLField()

    class Meta:
        ordering = ['-title']

class Population(models.Model):
    timestamp = models.DateTimeField(auto_now_add=True)
    server = models.ForeignKey(Server)
    players = models.PositiveIntegerField()

    class Meta:
        ordering = ['-timestamp', 'server']

