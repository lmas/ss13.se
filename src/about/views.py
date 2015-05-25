
from datetime import timedelta
from django.utils import timezone

from django.views import generic

from .utils import load_readme
from gameservers.models import Server, ServerHistory


class AboutView(generic.TemplateView):
    template_name = 'about/about.html'
    # The readme will only be loaded once at startup and then stored in mem.
    readme_md = load_readme()

    def get_context_data(self, **kwargs):
        context = super(AboutView, self).get_context_data(**kwargs)
        context['about_md'] = self.readme_md
        return context

class StatsView(generic.TemplateView):
    template_name = 'about/stats.html'
    delta_hour = timedelta(hours=1)
    delta_day = timedelta(days=1)

    def get_context_data(self, **kwargs):
        context = super(StatsView, self).get_context_data(**kwargs)

        now = timezone.now()
        last_update = now - self.delta_day
        total_players = 0
        servers_total = 0
        servers_with_players = 0
        servers_online = 0
        servers_warning = 0
        servers_offline = 0

        for server in Server.objects.all():
            servers_total += 1
            last_update = max(last_update, server.last_updated)

            if server.players_current > 0:
                total_players += server.players_current
                servers_with_players += 1

            delta = now - server.last_updated
            if delta < self.delta_hour:
                servers_online += 1
            elif delta < self.delta_day:
                servers_warning += 1
            else:
                servers_offline += 1

        context['last_update'] = last_update
        context['total_players'] = total_players
        context['servers_with_players'] = servers_with_players
        context['servers_total'] = servers_total
        context['servers_online'] = servers_online
        context['servers_warning'] = servers_warning
        context['servers_offline'] = servers_offline

        return context

