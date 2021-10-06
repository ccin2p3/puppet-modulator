# Metadata command

`pm` can help you with the common `metadata.json` modifications you may reproduce multiple times a day.

It can modify your `metadata.json` and also automatically commit the changes if you asked.

## Usage

```lang-none
❯ pm help metadata
Manipulate module metadata.json file

Usage:
  puppet-modulator metadata [command]

Available Commands:
  bump           
  puppet-version Interact with module required puppet version
  set-version    Set exact module version
  version        Get module version (current or next)

Flags:
  -g, --git-commit                       Commit changes to git
  -h, --help                             help for metadata
  -p, --keys-sort-commit-policy string   policy related to metadata keys sort commit. If pre-commit is used, then a dedicated commit will be created dedicated to metadata keys sorting. If no-pre-commit is used, metadata keys sorting will still occurs, but no dedicated commit will be created (default "pre-commit")
  -o, --output string                    Where to write metadata to. Defaults to modify metadata in-place

Global Flags:
      --config string   config file (default is $HOME/.puppet-modulator.yaml)
  -d, --debug           Enable debug

Use "puppet-modulator metadata [command] --help" for more information about a command.
```

## Automatically commit the modifications

In every `metadata` subcommand (`bump`, `set-version`), the flags `-g` / `--git-commit` and `-p` / `--keys-sort-commit-policy` control the _auto commit_ behavior.

If specified, the `-g` / `--git-commit` flag will enable the _auto commit_ behavior.

The flag `-p` / `--keys-sort-commit-policy` deserves a little more explanation.

### Rewrite the `metadata.json` file

<i class="fas fa-exclamation-circle"></i> Before digging in the _commit policy_, this is important to understand that `pm` **will rewrite your `metadata.json` file**.

To keep `pm` simple, this modification is made using the GO JSON standard library.

Your `metadata.json` file will be written following this spec:

* your `metadata.json` keys will be sorted alphabetically.
* your `metadata.json` will be indented with 2 spaces.

<br/>

The _alphabetical sort_ of the `metadata.json` file keys is not expected to change. This is the standard behavior of the GO JSON standard library.

On the other hand, the _hardcoded indent_ is quite easy to change and to customize to every personal preference. Feel free to [open an issue](https://github.com/ccin2p3/puppet-modulator/issues/new) to discuss.

<br/>

<i class="fas fa-info-circle"></i> The exact code responsible for this is:

```go
func (m MetadataJSON) WriteToWriter(w io.Writer, pretty bool) error {
  enc := json.NewEncoder(w)
  enc.SetIndent("", "  ")
  return enc.Encode(m.m)
}
```

### The metadata.json keys sort commit policy

Behind this barbarian name hides a very simple concept.

As said previously, as humans may edit the `metadata.json` file, it is likely that your keys order or your indent does not strictly matches the format `pm` will use.

The _keys sort commit policy_ controls how those changes will be commited to your `git` repository.

#### The `pre-commit` commit sort policy

The `pre-commit` commit sort policy is the default value.

Before applying the metadata modification you asked, `pm` will render the `metadata.json` in the format it expects, and commit those changes.
Only **after** those _aesthetic changes_ are commited, your real modification will be applied and commited on its own.

This allows to maintain a clear `diff` of the real modification you are introducing in your `git` history.

**Example**:

With an initial `metadata.json` file such as:

```json
{
  "author": "ccin2p3",
  "version": "4.0.0",
  "license": "CeCILL B",
  "name": "ccin2p3-modulator_test",
  "requirements": [
    {"name": "puppet","version_requirement": ">= 7.0.0 < 8.0.0"}
  ],
  "dependencies": [
    {
      "name": "puppetlabs/stdlib", "version_requirement": ">= 4.11.0 < 5.0.0"
    },
    {
      "name": "ccin2p3/etc_services",
      "version_requirement": ">= 2.0.0 < 3.0.0"
    }
  ],
  "summary": "Test module for puppet-modulator",
  "types": []
}
```

and a minor version changed commited with `pm metadata bump minor -g`, you'll see that two commits were created:

1. An _aesthetic changes_ commit only:

    ```diff
    diff --git a/metadata.json b/metadata.json
    index 684e0fd..c8f2d64 100644
    --- a/metadata.json
    +++ b/metadata.json
    @@ -1,20 +1,24 @@
    {
      "author": "ccin2p3",
    -  "version": "4.0.0",
    -  "license": "CeCILL B",
    -  "name": "ccin2p3-modulator_test",
    -  "requirements": [
    -    {"name": "puppet","version_requirement": ">= 7.0.0 < 8.0.0"}
    -  ],
      "dependencies": [
        {
    -      "name": "puppetlabs/stdlib", "version_requirement": ">= 4.11.0 < 5.0.0"
    +      "name": "puppetlabs/stdlib",
    +      "version_requirement": ">= 4.11.0 < 5.0.0"
        },
        {
          "name": "ccin2p3/etc_services",
          "version_requirement": ">= 2.0.0 < 3.0.0"
        }
      ],
    +  "license": "CeCILL B",
    +  "name": "ccin2p3-modulator_test",
    +  "requirements": [
    +    {
    +      "name": "puppet",
    +      "version_requirement": ">= 7.0.0 < 8.0.0"
    +    }
    +  ],
      "summary": "Test module for puppet-modulator",
    -  "types": []
    +  "types": [],
    +  "version": "4.0.0"
    }
    ```

    with the commit message `[meta] metadata.json automated modifications (pre-real modifications)`.

2. and the _real modification_ you were trying to make, the version bump, with a dedicated commit:

    ```diff
    diff --git a/metadata.json b/metadata.json
    index c8f2d64..f337595 100644
    --- a/metadata.json
    +++ b/metadata.json
    @@ -20,5 +20,5 @@
      ],
      "summary": "Test module for puppet-modulator",
      "types": [],
    -  "version": "4.0.0"
    +  "version": "4.1.0"
    }
    ```

    with the commit message `[meta] Bump version` (configurable with the `-m` / `--git-commit-msg` flag).

#### The `no-pre-commit` commit sort policy

As opposed to the `pre-commit` sort commit policy, the `no-pre-commit` will only create one commit that will include both the _aesthetic changes_ **and** the real modification you are making.

To reproduce the change above, with the same modification and the same initial `metadata.json` file, the command
`pm metadata bump minor -g --keys-sort-commit-policy no-pre-commit` will produce **one and only one** commit:

```diff
diff --git a/metadata.json b/metadata.json
index 684e0fd..f337595 100644
--- a/metadata.json
+++ b/metadata.json
@@ -1,20 +1,24 @@
 {
   "author": "ccin2p3",
-  "version": "4.0.0",
-  "license": "CeCILL B",
-  "name": "ccin2p3-modulator_test",
-  "requirements": [
-    {"name": "puppet","version_requirement": ">= 7.0.0 < 8.0.0"}
-  ],
   "dependencies": [
     {
-      "name": "puppetlabs/stdlib", "version_requirement": ">= 4.11.0 < 5.0.0"
+      "name": "puppetlabs/stdlib",
+      "version_requirement": ">= 4.11.0 < 5.0.0"
     },
     {
       "name": "ccin2p3/etc_services",
       "version_requirement": ">= 2.0.0 < 3.0.0"
     }
   ],
+  "license": "CeCILL B",
+  "name": "ccin2p3-modulator_test",
+  "requirements": [
+    {
+      "name": "puppet",
+      "version_requirement": ">= 7.0.0 < 8.0.0"
+    }
+  ],
   "summary": "Test module for puppet-modulator",
-  "types": []
+  "types": [],
+  "version": "4.1.0"
 }
```

with the commit message `[meta] Bump version`.

This makes it more difficult to find the _real modification_ we are introducing, but this saves you from an _aesthetic_ only commit.
