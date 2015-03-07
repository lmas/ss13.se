
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

        # Moving average for the last day
        # TODO: remove the hardcoded value
        tmp = [players for time, players in points[-96:]]
        context['daily_average'] = sum(tmp) / float(len(tmp))
        context['daily_min'] = min(tmp)
        context['daily_max'] = max(tmp)

        tmp = [players for time, players in points]
        context['total_average'] = sum(tmp) / float(len(tmp))
        context['total_min'] = min(tmp)
        context['total_max'] = max(tmp)
        return context

