# Puppet version requirements

Manipulate the `puppet` version requirements as described in the [modules metadata](https://puppet.com/docs/puppet/7/modules_metadata.html#metadata-version-requirement) official documentation.

## Usage

```lang-none
❯ pm metadata puppet-version -h
Interact with module required puppet version

Usage:
  puppet-modulator metadata puppet-version [command]

Available Commands:
  get         Get module required puppet version
  set         Set module required puppet version

[...]
```

## Get puppet version requirement

```lang-none
❯ pm metadata puppet-version get
>=4.10.0 <8.0.0
```

<i class="fas fa-info-circle"></i> The puppet version requirement constraint is validated and parsed.

## Set puppet version requirement

```lang-none
❯ pm metadata puppet-version set ">= 4.10.0 < 7.0.0"
```

<i class="fas fa-info-circle"></i> The puppet version requirement constraint passed as argument is validated and parsed.

<i class="fas fa-info-circle"></i> The flag `-g` / `--git-commit-msg` can be used to automatically commit the modifications.
