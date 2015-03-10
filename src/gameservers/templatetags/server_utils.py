
from django import template
from django.utils import timezone

register = template.Library()


@register.filter
def hours_since_now(timestamp):
    delta = timezone.now() - timestamp
    # Since a datetime.timedelta only stores the delta in days/seconds/ms
    # we have to convert the seconds to hours.
    return delta.total_seconds() / 3600.0

