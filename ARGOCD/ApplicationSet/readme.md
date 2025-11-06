1. So here create one manifest of Kind Application to maintain the Application manifest.
2. Then if you make any changes in the Kind application set it will monitor.

`🧠 So conceptually:

<!-- When you use a Git generator in an ApplicationSet like this:

generators:
  - git:
      repoURL: git@github.com:finternet-io/helmcharts.git
      targetRevision: v*
      directories:
        - path: charts/*


ArgoCD scans the repo and finds all matching directories.

For each matching directory (e.g. charts/yugabyte, charts/keycloak, etc.), it creates an object with metadata about that directory — such as:

Property	Meaning	Example Value
path	Full relative path of the directory	charts/yugabyte
path.basename	Just the last part of the path	yugabyte
path.filename	Full filename (if using file generator)	e.g. values.yaml
path.directory	Parent folder of the file/directory	e.g. charts -->