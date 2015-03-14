
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
        context['weekly_history'] = server.get_stats_history(days=7.5)
        context['averages_for_weekdays'] = server.get_averages_for_weekdays()
        return context

