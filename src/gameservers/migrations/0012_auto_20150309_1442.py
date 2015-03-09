# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0011_serverhistory'),
    ]

    operations = [
        migrations.AlterModelOptions(
            name='server',
            options={'ordering': ['-players_current', 'title']},
        ),
        migrations.RenameField(
            model_name='server',
            old_name='current_players',
            new_name='players_current',
        ),
        migrations.AddField(
            model_name='server',
            name='averages_for_weekdays',
            field=models.CommaSeparatedIntegerField(default=b'', max_length=50, editable=False),
            preserve_default=True,
        ),
        migrations.AddField(
            model_name='server',
            name='players_avg',
            field=models.PositiveIntegerField(default=0, editable=False),
            preserve_default=True,
        ),
        migrations.AddField(
            model_name='server',
            name='players_max',
            field=models.PositiveIntegerField(default=0, editable=False),
            preserve_default=True,
        ),
        migrations.AddField(
            model_name='server',
            name='players_min',
            field=models.PositiveIntegerField(default=0, editable=False),
            preserve_default=True,
        ),
    ]
