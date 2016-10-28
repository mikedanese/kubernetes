load("@io_bazel_rules_go//go:def.bzl", "go_env_attrs")

go_filetype = ["*.go"]

def _compute_make_variables(resolved_srcs, files_to_build):
  variables = {"SRCS": cmd_helper.join_paths(" ", resolved_srcs),
               "OUTS": cmd_helper.join_paths(" ", files_to_build)}
  if len(resolved_srcs) == 1:
    variables["<"] = list(resolved_srcs)[0].path
  if len(files_to_build) == 1:
    variables["@"] = list(files_to_build)[0].path
  return variables

def _go_sources_aspect_impl(target, ctx):
  transitive_sources = set(target.go_sources)
  for dep in ctx.rule.attr.deps:
    transitive_sources = transitive_sources | dep.transitive_sources
  return struct(transitive_sources = transitive_sources)

go_sources_aspect = aspect(
    attr_aspects = ["deps"],
    implementation = _go_sources_aspect_impl,
)

cmd_init_gopath = """
export GOROOT=$$(pwd)/%s/..
mkdir -p go/src/k8s.io/;
ln -s $$(pwd) go/src/k8s.io/kubernetes;
export GOPATH=$$(pwd)/go;
mkdir -p gengo/src/k8s.io/;
ln -s $$(pwd)/$(GENDIR) gengo/src/k8s.io/kubernetes;
export GOPATH=$${GOPATH}:$$(pwd)/gengo;
(
  cd go/src/k8s.io/kubernetes;
  %s;
)
"""

def _go_genrule_impl(ctx):
  all_srcs = set(ctx.files.go_src)
  label_dict = {}

  for dep in ctx.attr.go_deps:
    all_srcs = all_srcs | dep.transitive_sources

  for dep in ctx.attr.srcs:
    all_srcs = all_srcs | dep.files
    label_dict[dep.label] = dep.files

  resolved_inputs, argv, runfiles_manifests = ctx.resolve_command(
      command=cmd_init_gopath % (ctx.file.go_tool.dirname, ctx.attr.cmd.strip(chars=' \t\n\r')),
      attribute="cmd",
      expand_locations=True,
      make_variables=_compute_make_variables(all_srcs, set(ctx.outputs.outs)),
      tools=ctx.attr.tools,
      label_dict=label_dict
  )

  ctx.action(
      inputs = list(all_srcs) + resolved_inputs,
      outputs = ctx.outputs.outs,
      env = ctx.configuration.default_shell_env,
      command = argv,
      progress_message = "%s %s" % (ctx.attr.message, ctx),
      mnemonic = "GoGenrule",
  )

go_genrule = rule(
    attrs = go_env_attrs + {
        "srcs": attr.label_list(allow_files = True),
        "tools": attr.label_list(
            cfg = "host",
            allow_files = True,
        ),
        "outs": attr.output_list(mandatory = True),
        "cmd": attr.string(mandatory = True),
        "go_deps": attr.label_list(
            aspects = [go_sources_aspect],
        ),
        "message": attr.string(),
        "executable": attr.bool(default = False),
    },
    output_to_genfiles = True,
    implementation = _go_genrule_impl,
)
