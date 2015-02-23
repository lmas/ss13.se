
from django.shortcuts import render
from django.views import generic

from .models import Server, PlayerHistory

class ServerListView(generic.ListView):
    model = Server

class ServerDetailView(generic.DetailView):
    model = Server

    def get_context_data(self, **kwargs):
        context = super(ServerDetailView, self).get_context_data(**kwargs)
        server = context['server']
        history = PlayerHistory()
        points = history.get_points(server)
        context['player_history'] = points
        return context

