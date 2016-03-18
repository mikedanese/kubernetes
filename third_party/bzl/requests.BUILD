licenses(["notice"])

filegroup(
    name = "srcs",
    srcs = glob(["**"]),
    visibility = ["//third_party:__pkg__"],
)

py_library(
    name = "requests",
    srcs = [
        "requests/__init__.py",
        ":srcs",
    ],
    visibility = ["//visibility:public"],
)
