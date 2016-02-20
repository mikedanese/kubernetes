package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

genrule(
    name = "base_debian",
    srcs = [ "Dockerfile" ],
    outs = [ "base_debian.tar" ],
    local = 1,
    cmd = ""
    + "tag=\"gcr.io/google_containers:ansible_base\"; "
    + "cat $< | docker build -t $${tag} - && "
    + "docker save $${tag} >$@ && "
    + "docker rmi $${tag} "
    ,
)

docker_build(
    name = "base",
    repository = "gcr.io/google_containers",
    base = ":base_debian",
    files = [ "@jsonnet//:jsonnet" ],
    entrypoint = [ "/bin/bash", "-c" ],
)
