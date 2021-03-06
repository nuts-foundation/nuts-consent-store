Nuts consent store
==================

Consent storage for querying and checking active consent. This is the service space copy of the distributed consent.
Using the distributed records from the Corda vault is not possible due to the encrypted nature of those records.

.. image:: https://circleci.com/gh/nuts-foundation/nuts-consent-store.svg?style=svg
    :target: https://circleci.com/gh/nuts-foundation/nuts-consent-store
    :alt: Build Status

.. image:: https://readthedocs.org/projects/nuts-consent-store/badge/?version=latest
    :target: https://nuts-documentation.readthedocs.io/projects/nuts-consent-store/en/latest/?badge=latest
    :alt: Documentation Status

.. image:: https://codecov.io/gh/nuts-foundation/nuts-consent-store/branch/master/graph/badge.svg
    :target: https://codecov.io/gh/nuts-foundation/nuts-consent-store

.. image:: https://api.codeclimate.com/v1/badges/ddb963b417745047c472/maintainability
   :target: https://codeclimate.com/github/nuts-foundation/nuts-consent-store/maintainability
   :alt: Maintainability

.. include:: docs/pages/development/consent-store.rst
    :start-after: .. marker-for-readme

Configuration
*************

The following configuration parameters are available:

.. include:: README_options.rst

As with all other properties for nuts-go, they can be set through yaml:

.. sourcecode:: yaml

    cstore:
       connectionstring: :memory"

as commandline property

.. sourcecode:: shell

    ./nuts --cstore.connectionstring :memory:

Or by using environment variables

.. sourcecode:: shell

    NUTS_CSTORE_CONNECTIONSTRING=:memory: ./nuts
