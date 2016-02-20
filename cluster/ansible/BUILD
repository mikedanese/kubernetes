package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "ansible",
    files = glob([ "*" ]),
    package_dir = "/opt",
    mode = "0644",
)
