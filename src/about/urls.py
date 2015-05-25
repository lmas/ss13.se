
from django.conf.urls import patterns, url

from .views import AboutView, StatsView

urlpatterns = patterns('',
    url(r'^$', AboutView.as_view(), name='index'),
    url(r'^stats/$', StatsView.as_view(), name='stats'),
)
