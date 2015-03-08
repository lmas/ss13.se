
from django.shortcuts import render
from django.views import generic

from .models import Server, ServerHistory

class ServerListView(generic.ListView):
    model = Server

class ServerDetailView(generic.DetailView):
    model = Server

    def get_context_data(self, **kwargs):
        context = super(ServerDetailView, self).get_context_data(**kwargs)
        server = context['server']
        context['weekly_history'] = server.get_history_stats(days=7)

        avg, min, max = server.calc_player_stats(days=1)
        context['daily_average'] = avg
        context['daily_min'] = min
        context['daily_max'] = max

        avg, min, max = server.calc_player_stats(days=7)
        context['weekly_average'] = avg
        context['weekly_min'] = min
        context['weekly_max'] = max

        avg, min, max = server.calc_player_stats(days=31)
        context['monthly_average'] = avg
        context['monthly_min'] = min
        context['monthly_max'] = max
        return context

