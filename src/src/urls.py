from django.conf.urls import patterns, include, url
from django.contrib import admin
from django.views.generic.base import RedirectView


urlpatterns = patterns('',
    url(r'^admin/', include(admin.site.urls)),

    # HACK: redirect root index. Might have to do something about it soon..
    url(r'^$', RedirectView.as_view(url='/servers', permanent=False), name='index'),
    url(r'^servers/', include('gameservers.urls', namespace='gameservers')),
)
