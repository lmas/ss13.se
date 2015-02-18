from django.contrib import admin

from .models import Server, Population

class ServerAdmin(admin.ModelAdmin):
    list_display = ['title', 'site_url']

class PopulationAdmin(admin.ModelAdmin):
    list_display = ['timestamp', 'server', 'players']

admin.site.register(Server, ServerAdmin)
admin.site.register(Population, PopulationAdmin)

