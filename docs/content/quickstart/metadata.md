# Metadata operations overview

`pm` allows you to easily display the **current** version of your Puppet module.

```lang-none
❯ pm metadata version get
4.0.0
```

It can also helps you to set an arbitrary version in your `metadata.json` file:

```lang-none
❯ pm metadata set-version 42.0.0
```

Here is the _uncommited_ `diff` of the modification:

```diff
diff --git a/metadata.json b/metadata.json
--- a/metadata.json
+++ b/metadata.json
@@ -40,5 +40,5 @@
   "summary": "",
   "types": [],
-  "version": "4.0.0"
+  "version": "42.0.0"
 }
```

`pm` can also helps you **_bump_** your version in one command:

```lang-none
❯ pm metadata bump minor 
```

produces the following modification:

```diff
diff --git a/metadata.json b/metadata.json
--- a/metadata.json
+++ b/metadata.json
@@ -40,5 +40,5 @@
   "summary": "",
   "types": [],
-  "version": "4.0.0"
+  "version": "4.1.0"
 }
```

You can learn more in the [Metadata command](../metadata.md) dedicated page.