Nuts consent store
==================

Consent storage for querying and checking active consent. This is the service space copy of the distributed consent.
Using the distributed records from the Corda vault is not possible due to the encrypted nature of those records.

.. image:: https://travis-ci.org/nuts-foundation/nuts-consent-store.svg?branch=master
    :target: https://travis-ci.org/nuts-foundation/nuts-consent-store
    :alt: Build Status

.. image:: https://readthedocs.org/projects/nuts-consent-store/badge/?version=latest
    :target: https://nuts-documentation.readthedocs.io/projects/nuts-consent-store/en/latest/?badge=latest
    :alt: Documentation Status

.. image:: https://codecov.io/gh/nuts-foundation/nuts-consent-store/branch/master/graph/badge.svg
    :target: https://codecov.io/gh/nuts-foundation/nuts-consent-store

.. image:: https://api.codacy.com/project/badge/Grade/bbabffbece1a4ff4a34493f90078a66a
    :target: https://www.codacy.com/app/woutslakhorst/nuts-consent-store

The consent store is written in Go and should be part of nuts-go as an engine.

Dependencies
************

This projects is using go modules, so version > 1.12 is recommended. 1.10 would be a minimum. Currently Sqlite is used as database backend.

Running tests
*************

Tests can be run by executing

.. code-block:: shell

    go test ./...

Building
********

This project is part of https://github.com/nuts-foundation/nuts-go. If you do however would like a binary, just use ``go build``.

The client and server API is generated from the nuts-consent-store open-api spec:

.. code-block:: shell

    oapi-codegen -generate server,client,types -package api docs/_static/nuts-consent-store.yaml > api/generated.go


Generating mocks
----------------
Mocks used by other modules, generate with:

.. code-block:: shell

    mockgen -destination=mock/mock_client.go -package=mock -source=pkg/consent.go

Binary format migrations
------------------------

The database migrations are packaged with the binary by using the ``go-bindata`` package.

.. code-block:: shell

    NOT_IN_PROJECT $ go get -u github.com/go-bindata/go-bindata/...
    nuts-consent-store $ cd migrations && go-bindata -pkg migrations .

README
******

The readme is auto-generated from a template and uses the documentation to fill in the blanks.

.. code-block:: shell

    ./generate_readme.sh

This script uses ``rst_include`` which is installed as part of the dependencies for generating the documentation.

Documentation
*************

To generate the documentation, you'll need python3, sphinx and a bunch of other stuff. See :ref:`nuts-documentation-development-documentation`
The documentation can be build by running

.. code-block:: shell

    /docs $ make html

The resulting html will be available from ``docs/_build/html/index.html``

Configuration
*************

The following configuration parameters are available.

=====================================   ====================    ================================================================
Property                                Default                 Description
=====================================   ====================    ================================================================
nuts.cstore.connectionstring            :memory:                Sqlite connection string
nuts.cstore.mode                        server                  Server or client mode
nuts.cstore.address                     localhost:1323          Address of the server when in client mode
=====================================   ====================    ================================================================

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

