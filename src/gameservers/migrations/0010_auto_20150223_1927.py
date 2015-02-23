# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0009_server_last_updated'),
    ]

    operations = [
        migrations.AlterField(
            model_name='server',
            name='game_url',
            field=models.CharField(max_length=255),
            preserve_default=True,
        ),
    ]
