
local version = std.extVar("tag");

local name = "kubelet-ds";

local command = [
    "nsenter",
    "--target=1",
    "--mount",
    "--wd=.",
    "--",
    "kubelet",
    "--api-servers=https://kubernetes-master",
    "--enable-debugging-handlers=true",
    "--cloud-provider=gce",
    "--config=/etc/kubernetes/manifests",
    "--allow-privileged=True"
    "--v=4",
    "--cluster-dns=10.0.0.10",
    "--cluster-domain=cluster.local",
    # "--configure-cbr0=true",
    "--cgroup-root=/",
    "--system-container=/system",
];

local labels = {
    "tier": "node",
    "component": "kubelet",
    "version": version,
};

local objectMeta = {
    "metadata": {
        "name": name,
        "namespace": "kube-system",
        "labels": labels,
    },
};

local typeMeta = {
    "apiVersion": "extensions/v1beta1",
    "kind": "DaemonSet",
};

local podSpec = {
    "hostNetwork": true,
    "hostPID": true,
    "containers": [{
        "name": "kubelet",
        "image": "gcr.io/mikedanese-k8s/kubelet:%s" % [ version ],
        "args": command,
        "securityContext": {
            "privileged": true,
        },
        "imagePullPolicy": "Always",
    }]
};

local dsSpec = {
    "template": {
        "metadata": { "labels": labels },
        "spec": podSpec,
    },
};

objectMeta + typeMeta {
    "spec": dsSpec,
}
