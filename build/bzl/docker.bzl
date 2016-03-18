
def docker_pull(name, image, digest, registry="index.docker.io", repository="library"):
    out_path = "%s.tar" % name
    args = [
        "--registry="+registry,
        "--repository="+repository,
        "--image="+image,
        "--digest="+digest,
        "--out_path=$@",
    ]
    cmd = ""
    for part in ["$(location //build/bzl:docker_pull)"] + args:
        cmd += part
        cmd += " "

    native.genrule(
        name = name,
        outs = [out_path],
        cmd = cmd,
        local = 1,
        tools = [
            "//build/bzl:docker_pull",
        ],
    )
