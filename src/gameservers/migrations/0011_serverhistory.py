# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations
import django.utils.timezone


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0010_auto_20150223_1927'),
    ]

    operations = [
        migrations.CreateModel(
            name='ServerHistory',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('created', models.DateTimeField(default=django.utils.timezone.now)),
                ('players', models.PositiveIntegerField(default=0)),
                ('server', models.ForeignKey(to='gameservers.Server')),
            ],
            options={
                'ordering': ['-created', 'server'],
            },
            bases=(models.Model,),
        ),
    ]
