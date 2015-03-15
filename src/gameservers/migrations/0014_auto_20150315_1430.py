# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0013_auto_20150310_1944'),
    ]

    operations = [
        migrations.AlterModelOptions(
            name='server',
            options={'ordering': ['-players_current', '-last_updated', 'title']},
        ),
    ]
