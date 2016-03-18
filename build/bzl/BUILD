package(default_visibility = ["//visibility:public"])


py_binary(
    name = "docker_pull",
    srcs = ["docker_pull.py"],
    deps = [
        "@bazel_tools//third_party/py/gflags",
        "@requests//:requests",
    ],
)
