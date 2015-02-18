# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('gameservers', '0005_auto_20150218_2029'),
    ]

    operations = [
        migrations.CreateModel(
            name='PopulationHistory',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('timestamp', models.DateTimeField(auto_now_add=True)),
                ('players', models.PositiveIntegerField()),
                ('server', models.ForeignKey(to='gameservers.Server')),
            ],
            options={
                'ordering': ['timestamp', 'server'],
            },
            bases=(models.Model,),
        ),
        migrations.RemoveField(
            model_name='population',
            name='server',
        ),
        migrations.DeleteModel(
            name='Population',
        ),
    ]
