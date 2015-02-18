from django.contrib import admin

from .models import Server, PopulationHistory

class ServerAdmin(admin.ModelAdmin):
    list_display = ['title', 'site_url']
    search_fields = ['title']

class PopulationHistoryAdmin(admin.ModelAdmin):
    list_display = ['timestamp', 'players', 'server']
    list_filter = ['timestamp']
    search_fields = ['server__title']

admin.site.register(Server, ServerAdmin)
admin.site.register(PopulationHistory, PopulationHistoryAdmin)

