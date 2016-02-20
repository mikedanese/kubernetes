package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

docker_build(
    name = "k8s_installer",
    repository = "gcr.io/google_containers",
    base = "//cluster/ansible/base",
    tars = [ "//cluster/ansible", "//:release" ],
    entrypoint = [ "/bin/bash", "-c" ],
)
