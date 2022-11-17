#!/bin/sh
set -u
# Source the vault database secrets from deployment if present before starting the service
if [ -f /vault/secrets/database-config ]; then
  source /vault/secrets/database-config
else
  echo "Warning: Vault secrets file /vault/secrets/database-config not found" >&2
fi

# execute passed-in parameters
$@

# capture main command exit code
mainExit=$?

# Workaround for https://github.com/istio/istio/issues/6324 if a job
# This is never reached if the above doesn't terminate (i.e. not a job)
curl -sfI -X POST http://127.0.0.1:15020/quitquitquit

# exit with main job's exit code
exit $mainExit
