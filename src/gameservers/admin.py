from django.contrib import admin

from .models import Server

class ServerAdmin(admin.ModelAdmin):
    list_display = ['title', 'site_url']
    search_fields = ['title']

admin.site.register(Server, ServerAdmin)

