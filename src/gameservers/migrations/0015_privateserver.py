# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0014_auto_20150315_1430'),
    ]

    operations = [
        migrations.CreateModel(
            name='PrivateServer',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('title', models.CharField(max_length=255)),
                ('site_url', models.URLField(blank=True)),
                ('host', models.CharField(max_length=255)),
                ('port', models.PositiveIntegerField()),
                ('active', models.BooleanField(default=False)),
            ],
            options={
                'ordering': ['-active', 'title'],
            },
        ),
    ]
