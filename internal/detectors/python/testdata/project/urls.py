# TEST: ignores triple double quotes
"""
ignore
"""

CONFIG = {
    # TEST: ignores string keys with double quotes
    # TEST: detects string literals from dict values
    "ignore.domain.com": "string-key.example.com",
    # TEST: ignores string keys with single quotes
    'ignore.domain.com': 42
}

def f():
    # TEST: ignores triple single quotes
    '''ignore.domain.com'''

    user_id = 42

    # TEST: detects string literals in expressions
    hostname = "variable.example.com"
    # TEST: detects concatenated strings
    hostname = "concat" + ".example.com"
    # TEST: detects strings with interpolation
    hostname = f'a.{user_id}.interpolation.example.com'
    # TEST: detects raw strings
    hostname = r"raw.example.com"
    # TEST: ignores indices
    hostname = CONFIG["ignore.domain.com"]

    return hostname


# TEST: detects concatenated strings as a single value
concat_string = (
    "part1."
    "example.com"
)

@connect_on_app_finalize
def add_backend_cleanup_task(app):
    @app.task(name='celery.backend_cleanup', shared=False, lazy=False)
    def backend_cleanup():
        app.backend.cleanup()
    return backend_cleanup


@DictProperty('environ', 'ignore.domain.com', read_only=True)
def app(self):
    """ Bottle application handling this request. """
    raise RuntimeError('This request is not connected to an application.')
