# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Population',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('timestamp', models.DateTimeField(auto_now_add=True)),
                ('players', models.PositiveIntegerField()),
            ],
            options={
                'ordering': ['-timestamp', 'server'],
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Server',
            fields=[
                ('id', models.AutoField(verbose_name='ID', serialize=False, auto_created=True, primary_key=True)),
                ('title', models.CharField(max_length=255)),
                ('game_url', models.URLField()),
                ('site_url', models.URLField()),
            ],
            options={
                'ordering': ['-title'],
            },
            bases=(models.Model,),
        ),
        migrations.AddField(
            model_name='population',
            name='server',
            field=models.ForeignKey(to='gameservers.Server'),
            preserve_default=True,
        ),
    ]
