# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0003_auto_20150218_1652'),
    ]

    operations = [
        migrations.AlterModelOptions(
            name='population',
            options={'ordering': ['timestamp', 'server']},
        ),
    ]
