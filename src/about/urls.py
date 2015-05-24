
from django.conf.urls import patterns, url

from .views import AboutView

urlpatterns = patterns('',
    url(r'^$', AboutView.as_view(), name='index'),
)
