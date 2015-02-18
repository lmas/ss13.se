
from django.conf.urls import patterns, url

from .views import ServerListView, ServerDetailView

urlpatterns = patterns('',
    url(r'^$', ServerListView.as_view(), name='index'),

    url(r'^(?P<pk>\d+)/$', ServerDetailView.as_view(), name='detail'),
    url(r'^(?P<pk>\d+)/(?P<slug>.+)/$', ServerDetailView.as_view(), name='detail'),
)
