package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "release",
    files = [
        "_output/release-stage/server/linux-amd64/kubernetes/server/bin/kubelet",
        "_output/release-stage/server/linux-amd64/kubernetes/server/bin/kubectl",
    ],
    package_dir = "/usr/local/bin",
    mode = "0777",
)
