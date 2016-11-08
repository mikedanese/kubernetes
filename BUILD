package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load("@io_bazel_rules_go//go:def.bzl", "go_prefix")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

go_prefix("k8s.io/kubernetes")

pkg_tar(
    name = "kubernetes-sources",
    files = glob(
        ["**"],
        exclude = [
            ".git/**/*",
            "_output/**/*",
            "bazel-*/**/*",
        ],
    ),
    strip_prefix = ".",
)
