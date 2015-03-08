
from datetime import timedelta

from django.shortcuts import render
from django.views import generic
from django.utils import timezone

from .models import Server, ServerHistory

class ServerListView(generic.ListView):
    model = Server

class ServerDetailView(generic.DetailView):
    model = Server

    def get_context_data(self, **kwargs):
        context = super(ServerDetailView, self).get_context_data(**kwargs)
        server = context['server']

        weekly_history = ServerHistory.objects.filter(
            server=server,
            created__gte=timezone.now() - timedelta(days=7),
        )
        context['weekly_history'] = weekly_history

        # Moving average for the last day
        tmp = [tmp.players for tmp in ServerHistory.objects.filter(
            server=server,
            created__gte=timezone.now() - timedelta(days=1))]
        context['daily_average'] = sum(tmp) / float(len(tmp))
        context['daily_min'] = min(tmp)
        context['daily_max'] = max(tmp)

        tmp = [tmp.players for tmp in weekly_history]
        context['weekly_average'] = sum(tmp) / float(len(tmp))
        context['weekly_min'] = min(tmp)
        context['weekly_max'] = max(tmp)
        return context

