<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/godel-okgo-asset-outparamcheck"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

godel-okgo-asset-outparamcheck
==============================
godel-okgo-asset-outparamcheck is an asset for the gödel [okgo plugin](https://github.com/palantir/okgo). It provides
the functionality of the [outparamcheck](https://github.com/palantir/outparamcheck) check.

This check verifies that output parameters are properly passed in to functions as pointers.

Configuration
-------------
It is possible to specify functions that should be checked using output parameters using configuration. The check
configuration has an `out-param-fns` key that stores a map from the fully qualified function name to a slice of argument
indices that indicate the indices of the function parameters that should be considered output parameters.

Here is an example `check.yml` configuration that specifies that the second (index 1) argument of the "Load" function in
the "github.com/org/repo/config" package is an output parameter:

```yaml
checks:
  outparamcheck:
    config:
      out-param-fns:
        "github.com/org/repo/config.Load": [1]
```
