
from django import template
from django.utils import timezone

register = template.Library()


@register.filter
def hours_since_now(timestamp):
    delta = timezone.now() - timestamp
    return delta.total_seconds() / 3600.0

