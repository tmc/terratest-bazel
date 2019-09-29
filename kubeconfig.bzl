_doc = """
kubeconfig copies kubernetes configuration into the bazel workspace.
"""

def _kubeconfig_impl(ctx):
    quiet = True
    er = ctx.execute(["kubectl", "config", "view", "--raw"], timeout = 5, quiet = quiet)
    if er.return_code != 0:
        print("issue running kubectl:", er.stderr)
    er = ctx.execute(["kind","get","kubeconfig"], timeout = 5, quiet = quiet)
    ctx.file("kubeconfig.yaml", content = er.stdout)
    ctx.file("BUILD", content = """exports_files(glob(['*']))""")

kubeconfig_rule = repository_rule(
    _kubeconfig_impl,
    local = True,
    doc = _doc,
    environ = ["KUBECONFIG"],
)

def kubeconfig():
    kubeconfig_rule(name = "kubeconfig")
