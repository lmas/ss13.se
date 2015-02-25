
from fabric.api import env, task, run, warn_only, cd, prefix, sudo
from fabric.contrib.files import exists

from contextlib import contextmanager as _contextmanager

env.use_ssh_config = True

env.user = 'lmas'
env.project = 'ss13_hub'

env.path = '/opt/{}'.format(env.project)
env.repo_path = '{}/{}'.format(env.path, env.project)
env.venv_path = '{}/venv'.format(env.path)
env.activate = 'source {}/bin/activate'.format(env.venv_path)

@_contextmanager
def virtualenv():
    with cd(env.repo_path):
        with prefix(env.activate):
            yield

################################################################################

@task
def init_base_dir():
    if exists(env.path):
        return
    sudo('mkdir {}'.format(env.path))
    sudo('chown {0}:{0} {1}'.format(env.user, env.path))

@task
def update_repo():
    '''Deploy the app to a server.'''
    if exists(env.repo_path):
        with cd(env.repo_path):
            run('git pull')
    else:
        run('git clone git://github.com/lmas/{}.git {}'.format(env.project, env.repo_path))

@task
def init_venv():
    if exists(env.venv_path):
        return
    with cd(env.path):
        run('virtualenv {}'.format(env.venv_path))

@task
def update_venv():
    with virtualenv():
        run('pip install --upgrade -r {}/requirements.txt'.format(env.repo_path))

@task
def collect_static():
    # TODO: needs to grab settings before this
    with virtualenv():
        run('python src/manage.py collectstatic --noinput')

@task
def migrate():
    with virtualenv():
        run('python src/manage.py migrate')
        # TODO: create a superuser?

@task
def init_supervisor():
    if exists('/etc/supervisor/conf.d/'):
        sudo('cp {}/supervisor.conf /etc/supervisor/conf.d/{}.conf'.format(
            env.repo_path, env.project))
        sudo('supervisorctl reread')
        sudo('supervisorctl add {}'.format(env.project))

@task
def restart():
    sudo('supervisorctl restart {}'.format(env.project))

@task
def deploy():
    '''Initialize some directories and deploy the app on a server.'''
    init_base_dir()
    update_repo()
    init_venv()
    update_venv()
    collect_static()
    migrate()
    init_supervisor()
    restart()

@task
def quick_update():
    '''Update to latest version of the app and restart.'''
    update_repo()
    collect_static()
    migrate()
    restart()

