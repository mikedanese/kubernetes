load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "h2stream.go",
        "pipe.go",
        "v5_codec.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/util/h2stream",
    importpath = "k8s.io/apimachinery/pkg/util/h2stream",
    visibility = ["//visibility:public"],
    deps = ["//vendor/golang.org/x/net/http2:go_default_library"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["h2stream_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//staging/src/k8s.io/client-go/util/cert:go_default_library",
        "//vendor/golang.org/x/net/http2:go_default_library",
    ],
)
