#! /bin/bash

set -o errexit
set -o nounset
set -o pipefail

declare -r HOST_ROOT="/host_root"

echo '{}' > componentconfig.json

cat <<EOF | jsonnet \
  --code-var ohai="$(ohai)" \
  --var host_root="${HOST_ROOT}" \
  --output-file in.json -

local ohai = std.extVar("ohai");
local host_root = std.extVar("host_root");

{
    host_root: host_root,
    componentconfig: import "componentconfig.json",
    playbooks: import "books.json",
    ohai: ohai,
}
EOF

cat <<EOF > hosts
[node]
${HOST_ROOT}

[master]
${HOST_ROOT}

[minion]
${HOST_ROOT}
EOF

playbooks=()
for playbook in $(jq -c '.playbooks[]' in.json); do
  name=$(echo "${playbook}" | jq -r '.name')
  echo "${name}"

  echo "${playbook}" | jq '.params' > "roles/${name}/params.json"

  playbooks+=("roles/${name}/run.yml")
done

ansible-playbook \
  --inventory-file=hosts \
  --connection=chroot \
  ${playbooks[@]}
