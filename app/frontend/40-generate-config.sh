#!/bin/sh

cat > /usr/share/nginx/html/config.js <<EOF
window.__APP_CONFIG__ = {
  apiBaseUrl: "${API_BASE_URL}",
  webtransportEndpoint: "${WEBTRANSPORT_ENDPOINT}",
  useCertHash: ${WEBTRANSPORT_USE_CERT_HASH:-false},
  certHash: "${WEBTRANSPORT_CERT_HASH}"
};
EOF