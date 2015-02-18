from django.contrib import admin

from .models import Server, PopulationHistory

class ServerAdmin(admin.ModelAdmin):
    list_display = ['title', 'site_url']

class PopulationHistoryAdmin(admin.ModelAdmin):
    list_display = ['timestamp', 'server', 'players']

admin.site.register(Server, ServerAdmin)
admin.site.register(PopulationHistory, PopulationHistoryAdmin)

