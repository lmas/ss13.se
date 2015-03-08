
from django.shortcuts import render
from django.views import generic

from .models import Server

class ServerListView(generic.ListView):
    model = Server

class ServerDetailView(generic.DetailView):
    model = Server

    def get_context_data(self, **kwargs):
        context = super(ServerDetailView, self).get_context_data(**kwargs)
        server = context['server']
        context['weekly_history'] = server.get_history_stats(days=7)

        stats = server.calc_player_stats(days=1)
        context['daily_average'] = stats['players__avg']
        context['daily_min'] = stats['players__min']
        context['daily_max'] = stats['players__max']

        stats = server.calc_player_stats(days=7)
        context['weekly_average'] = stats['players__avg']
        context['weekly_min'] = stats['players__min']
        context['weekly_max'] = stats['players__max']

        stats = server.calc_player_stats(days=31)
        context['monthly_average'] = stats['players__avg']
        context['monthly_min'] = stats['players__min']
        context['monthly_max'] = stats['players__max']

        context['weekday_averages'] = server.weekday_averages()
        return context

