# Frequently Asked Questions

## `metadata.json` structure is rewritten after change

After a `metadata.json` modification, the `metadata.json` file structure _may_  change.

This is simply due to the fact that `pm` is using the standard GO JSON library.

When the modified JSON content is written to the `metadata.json` file, keys will automatically be
sorted in alphabetical order. [`puppet-blacksmith` does have the exact same behavior](https://github.com/voxpupuli/puppet-blacksmith/issues/96).

This is expected behavior and there is no way to change that.

The only thing you can control is how this will _structure modification_ will be commited to your repository. The [`metadata.json` keys sort commit policy documentation](metadata.md#the-metadatajson-keys-sort-commit-policy) will explain in detail this behavior.

## I'd like to use the _version auto-guess_ feature in `flow` and still specify a _base reference branch_

If you want to specify a `base reference` branch and still use the _version auto-guess_ feature, you can use `""` (empty string) for the version, or `?` (question mark).

```lang-none
❯ pm flow release start "?" support/1.x.x
Switched to a new branch 'release/1.3.0'

Summary of actions:
- A new branch 'release/1.3.0' was created, based on 'support/1.x.x'
- metadata.json version was set to 1.3.0 and automatically commited for you
- You are now on branch 'release/1.3.0'

[...]
```

## I'd like to report a bug / suggest a change / contribute

The public github repository is here for that.

Feel free to [open an issue](https://github.com/ccin2p3/puppet-modulator/issues/new) or [create a pull request](https://github.com/ccin2p3/puppet-modulator/pulls).