package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    visibility = ["//visibility:private"],
    deps = ["//cmd/gals/cmd:go_default_library"],
)

go_binary(
    name = "gals",
    library = ":go_default_library",
    linkstamp = "istio.io/galley/pkg/version",
    visibility = ["//visibility:public"],
)
