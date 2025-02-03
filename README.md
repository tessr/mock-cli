# mock-cli

This repo is a place where I can prototype and play with CLI tools. These tools are **mockups**: They do not do anything, but they let me play with the UX.

## netchef 

### Installation 

```
go install github.com/tessr/mock-cli/cmd/netchef
```

### Usage 
Generate the required files for devnet with reasonable defaults and randomly generated name / chain ID:

```
$ netchef generate
```

Deploy the devnet, using the files created through `netchef generate`:

```
$ netchef deploy
```
