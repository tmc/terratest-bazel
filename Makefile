
update-build-files:
	bazel run //:gazelle

update-go-deps:
	bazel run //:gazelle -- update-repos -from_file=go.mod -prune -to_macro=third_party/godeps.bzl%go_dependencies
