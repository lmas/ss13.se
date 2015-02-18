# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0004_auto_20150218_1915'),
    ]

    operations = [
        migrations.AlterModelOptions(
            name='server',
            options={'ordering': ['-current_players', 'title']},
        ),
        migrations.AddField(
            model_name='server',
            name='current_players',
            field=models.PositiveIntegerField(default=0),
            preserve_default=True,
        ),
    ]
