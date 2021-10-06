# Module version operations

## Bump patch version

```lang-none
❯ pm metadata bump patch
```

produces the following `diff`:

```diff
diff --git a/metadata.json b/metadata.json
index f337595..152cf97 100644
--- a/metadata.json
+++ b/metadata.json
@@ -20,5 +20,5 @@
   ],
   "summary": "Test module for puppet-modulator",
   "types": [],
-  "version": "4.1.0"
+  "version": "4.1.1"
 }
```

## Bump minor version

```lang-none
❯ pm metadata bump minor
```

produces the following `diff`:

```diff
diff --git a/metadata.json b/metadata.json
index f337595..152cf97 100644
--- a/metadata.json
+++ b/metadata.json
@@ -20,5 +20,5 @@
   ],
   "summary": "Test module for puppet-modulator",
   "types": [],
-  "version": "4.1.0"
+  "version": "4.2.0"
 }
```

## Bump major version

```lang-none
❯ pm metadata bump minor
```

produces the following `diff`:

```diff
diff --git a/metadata.json b/metadata.json
index f337595..152cf97 100644
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

## Arbitrary set a version

```lang-none
❯ pm metadata set-version 42.0.0
```

produces the following `diff`:

```diff
diff --git a/metadata.json b/metadata.json
index f337595..763de0c 100644
--- a/metadata.json
+++ b/metadata.json
@@ -20,5 +20,5 @@
   ],
   "summary": "Test module for puppet-modulator",
   "types": [],
-  "version": "4.1.0"
+  "version": "42.0.0"
 }
```

## Get module version

```lang-none
❯ pm metadata version get
4.1.0
```
