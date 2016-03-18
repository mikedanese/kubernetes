package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("//build/bzl:docker.bzl", "docker_pull")

pkg_tar(
    name = "ansible",
    files = glob([ "*" ]),
    package_dir = "/opt",
    mode = "0644",
)

docker_pull(
    name = "busybox",
    registry = "gcr.io",
    repository = "google_containers",
    image = "busybox",
    digest = "sha256:4bdd623e848417d96127e16037743f0cd8b528c026e9175e22a84f639eca58ff",
)
