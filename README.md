# git-credential-token
Git credential helper to put access tokens behind a password when two-factor is enabled on GitHub.

TODO: Update this on how to use since the git plugins setup isn't very obvious.

[credential "https://king-jam@github.com"] <- username not respected....
        helper = crypt-store -file ~/.git-credential-crypt-store
