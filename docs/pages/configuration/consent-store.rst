.. _nuts-consent-store-configuration:

Nuts consent store configuration
################################

.. marker-for-readme

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