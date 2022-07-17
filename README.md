# go-git-credential-split

It's an alternative credential helper for git. It allows you to store and use credentials in multiple files.

I'm building this project to give myself more convenience and additional security. If you're working on multiple 
repositories with different accounts/credentials, in my case for my hobby and professional, you can store credentials 
in more than one file. I keep projects for my hobby and professional on separate encrypted volumes along with 
the respective credentials. If you have a similar use case or any other use case you think the tool may work for you, 
feel free to use it.

## Build & Install

    go build -o git-credential-split

you may need to move/copy it to your `PATH`

    cp git-credential-split /usr/local/bin/

### Configure git to use credential split

    git config credential.helper split
    git config credential.usehttppath true

## Warning

> :warning: This project is still in development. While it may satisfy my requirement for now, it still lacks of core
> functionality.