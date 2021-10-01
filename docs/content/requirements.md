# Requirements

## Working directory

`pm` expects you to work from the _root directory_ of your Puppet module. In other words, the directory your `metadata.json` file lives in.


<i class="fas fa-exclamation-circle"></i> Multiple `pm` `metadata` subcommands expect the working directory to be a `git` managed repository.

<i class="fas fa-exclamation-circle"></i> `pm` `flow` subcommand also expects the working directory to be a [git-flow](https://nvie.com/posts/a-successful-git-branching-model/) managed repository.