from django.contrib import admin

from .models import PrivateServer, Server

class PrivateServerAdmin(admin.ModelAdmin):
    list_display = ['title', 'site_url', 'active']
    search_fields = ['title']

class ServerAdmin(admin.ModelAdmin):
    list_display = ['title', 'site_url']
    search_fields = ['title']

admin.site.register(PrivateServer, PrivateServerAdmin)
admin.site.register(Server, ServerAdmin)

