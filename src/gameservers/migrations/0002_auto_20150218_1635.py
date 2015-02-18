# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0001_initial'),
    ]

    operations = [
        migrations.AlterField(
            model_name='server',
            name='site_url',
            field=models.URLField(blank=True),
            preserve_default=True,
        ),
    ]
