.. _nuts-consent-store-api:

Nuts consent store API
======================

.. raw:: html

    <div id="swagger-ui"></div>

    <script src='../../_static/js/swagger-ui-bundle-3.18.3.js' type='text/javascript'></script>
    <script src='../../_static/js/swagger-ui-standalone-preset-3.18.3.js' type='text/javascript'></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                "dom_id": "#swagger-ui",
                urls: [
                    {url: "../../_static/nuts-consent-store.yaml", name: "consent-store"},
                    ],
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                layout: "StandaloneLayout"
            });

            window.ui = ui
        }

    </script>