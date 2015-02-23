# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0007_auto_20150218_2107'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='populationhistory',
            name='server',
        ),
        migrations.DeleteModel(
            name='PopulationHistory',
        ),
    ]
