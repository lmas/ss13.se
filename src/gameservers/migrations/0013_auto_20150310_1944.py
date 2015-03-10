# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations
import django.utils.timezone


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0012_auto_20150309_1442'),
    ]

    operations = [
        migrations.AlterField(
            model_name='server',
            name='last_updated',
            field=models.DateTimeField(default=django.utils.timezone.now, editable=False),
            preserve_default=True,
        ),
    ]
