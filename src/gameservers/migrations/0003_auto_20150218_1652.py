# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0002_auto_20150218_1635'),
    ]

    operations = [
        migrations.AlterModelOptions(
            name='server',
            options={'ordering': ['title']},
        ),
    ]
