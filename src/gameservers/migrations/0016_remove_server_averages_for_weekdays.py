# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0015_privateserver'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='server',
            name='averages_for_weekdays',
        ),
    ]
