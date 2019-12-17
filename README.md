# git-credential-token
Git credential helper to put access tokens behind a password when two-factor is enabled on GitHub.

Can be run on Linux using zenity for dialog boxes.

Add the following to your git config file:

``` git
[credential]
        helper = crypt-store
```
