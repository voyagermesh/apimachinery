# Versioning Policy

There are 2 parts to versioning policy:
 - Operator version: Voyager __does not follow semver__, rather the _major_ version of operator points to the
Kubernetes [client-go](https://github.com/kubernetes/client-go#branches-and-tags) version. You can verify this
from the `glide.yaml` file. This means there might be breaking changes between point releases of the operator.
This generally manifests as changed annotation keys or their meaning.
Please always check the release notes for upgrade instructions.
 - TPR version: appscode.com/v1beta1 is considered in beta. This means any changes to the YAML format will be backward
compatible among different versions of the operator.
