
import hashlib

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
        context['graph_file'] = hashlib.sha256(server.title).hexdigest()
        return context

