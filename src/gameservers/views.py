from django.shortcuts import render
from django.views import generic

from .models import Server, PopulationHistory

class ServerListView(generic.ListView):
    model = Server

class ServerDetailView(generic.DetailView):
    model = Server

    def get_context_data(self, **kwargs):
        context = super(ServerDetailView, self).get_context_data(**kwargs)
        server = context['server']
        # HACK: 24 hours for the last 3 days, might want to change this
        context['population'] = PopulationHistory.objects.filter(server=server)[:3*24]
        return context

