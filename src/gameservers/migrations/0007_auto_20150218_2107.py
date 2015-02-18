# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0006_auto_20150218_2033'),
    ]

    operations = [
        migrations.AlterField(
            model_name='server',
            name='current_players',
            field=models.PositiveIntegerField(default=0, editable=False),
            preserve_default=True,
        ),
    ]
