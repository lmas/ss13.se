
import codecs
from os import path

import markdown as md

from django.conf import settings


def markdown(text):
    # In case I would ever switch to some other markdown lib.
    return md.markdown(text)

def load_readme():
    '''Read the contents of the project README and return it as markdown html.'''
    readme_path = path.abspath(path.join(settings.BASE_DIR, '..', 'README.md'))
    with codecs.open(readme_path, mode='r', encoding='utf-8') as f:
        tmp = f.read()
    return markdown(tmp)

