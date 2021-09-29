# git flow operations

## Why integrate such `git-flow` wrapper ?

At IN2P3-CC, we're using `git-flow` to manage all our modules.

The `git-flow` `release` and `hotfix` operations include a `metadata.json` interaction to change the version.

We also always need to check what version we're in to deduce the _next version_ of our module.

`pm` fits exactly here and helps you:

  * set your version in the `metadata.json` file and commit the change
  * in _auto mode_, guess the next `hotfix` or `release` version for you, then edit and commit

## Usage

### Release

```lang-none
❯ pm help flow release 
A git-flow high-level wrapper for releases

Usage:
  puppet-modulator flow release [command]

Available Commands:
  finish      A git-flow high-level wrapper to finish releases
  start       A git-flow high-level wrapper to start releases

Flags:
  -h, --help   help for release

Global Flags:
      --config string   config file (default is $HOME/.puppet-modulator.yaml)
  -d, --debug           Enable debug

Use "puppet-modulator flow release [command] --help" for more information about a command.
```

### Hotfix

```lang-none
❯ pm help flow hotfix 
A git-flow high-level wrapper for hotfixes

Usage:
  puppet-modulator flow hotfix [command]

Available Commands:
  finish      A git-flow high-level wrapper to finish hotfixes
  start       A git-flow high-level wrapper to start hotfixes

Flags:
  -h, --help   help for hotfix

Global Flags:
      --config string   config file (default is $HOME/.puppet-modulator.yaml)
  -d, --debug           Enable debug

Use "puppet-modulator flow hotfix [command] --help" for more information about a command.
```

## Release and Hotfix version auto-guess

If you do not specify a version in your `pm flow hotfix start` or `pm flow release start` command, `pm` will admit that you're trying to work quickly and use the most common version bump logic for those operations:

  * for a `release`, it will _increment the minor version_.
  * for a `hotfix`, it will _increment the patch version_.

**Important**:

  * If you don't want to use the _version auto-guess_ feature, you'll have to explicitly specify a version on command-line.
  * If you want to specify a `base reference` branch and still use the _auto-guess_ feature, you can use `""` (empty string) for the version, or `?` (question mark).

## Release

### Start a release

```lang-none
❯ pm flow release start -h
A git-flow high-level wrapper to start releases

Usage:
  puppet-modulator flow release start [version] [base-ref] [flags]

Flags:
  -h, --help   help for start

Global Flags:
      --config string   config file (default is $HOME/.puppet-modulator.yaml)
  -d, --debug           Enable debug
```

**Example**:

We're dealing with a module in version `4.1.0` and we want to start the release `5.0.0`:

```lang-none
❯ pm flow release start 5.0.0
Switched to a new branch 'release/5.0.0'

Summary of actions:
- A new branch 'release/5.0.0' was created, based on 'develop'
- metadata.json version was set to 5.0.0 and automatically commited for you
- You are now on branch 'release/5.0.0'

Follow-up actions:
- Start committing last-minute fixes in preparing your release
- When done, run:

        puppet-modulator gflow release finish

        or

        puppet-modulator gflow release finish -p -q -m "MESSAGE"
```

will set you on the branch `release/5.0.0` and has automatically created the commit:

```diff
diff --git a/metadata.json b/metadata.json
index f337595..c99e703 100644
--- a/metadata.json
+++ b/metadata.json
@@ -20,5 +20,5 @@
   ],
   "summary": "Test module for puppet-modulator",
   "types": [],
-  "version": "4.1.0"
+  "version": "5.0.0"
 }
```

### Finish a release

This is a simple wrapper on top of the `git flow release finish` command.

```lang-none
❯ pm flow release finish -h
A git-flow high-level wrapper to finish releases

Usage:
  puppet-modulator flow release finish [flags]

Flags:
  -h, --help                 help for finish
  -m, --message string       Use the given tag message
  -f, --messagefile string   Use the contents of the given file as tag message
  -q, --no-prompt            No prompt for editor (set GIT_MERGE_AUTOEDIT=no)
  -p, --push                 Push to origin after performing finish

Global Flags:
      --config string   config file (default is $HOME/.puppet-modulator.yaml)
  -d, --debug           Enable debug
```

The main addition is the `-q` / `--no-prompt` flag that allows you to bypass the _editor prompt_ that `git flow` does.

## Hotfix

`hotfix` and `release` are very similar. All the [release](#release) documentation applies here.