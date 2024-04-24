# tinybird-cli

A simple CLI tool to manage some bulk operations for TinyBird API (at this moment, only one operation [PUT](https://www.tinybird.co/docs/api-reference/token-api#put--v0-tokens-(.+)) is supported).
The app is built using the [cobra](https://github.com/spf13/cobra) CLI framework, and it has the following structure:

```test
├── cmd
│   ├── root.go
│   └── tokens
│       ├── tokens.go
│       └── actions        
│           └── put.go
└── main.go
```

### Flags:

| Flag         | Short Form | Description                                                                                                                                |
| ------------ | --- |--------------------------------------------------------------------------------------------------------------------------------------------|
| --admin-token| -a  | Workspace admin token                                                                                                                      |
| --file       | -f  | Path to the file containing tokens                                                                                                         |
| --decrypt-key| -d  | Encryption key (to be used exclusively when the source file contains encrypted tokens)                                                     |
| --scope      | -s  | This flag is used to provide a list of scopes ([for tokens API](https://www.tinybird.co/docs/api-reference/token-api#put--v0-tokens-(.+))) |

## Usage
According to directory structure app args looks as follows
```
tinybird-cli [entity] [action] [flags]
```

Since the goal was to create a simple wrapper over the TB API with minimal dependencies, the data source (specifically, a list of tokens) is passed as a file path (txt or csv without a header)
The file can include either raw tokens or tokens encrypted using AES256ECB encryption.

Example file content (each token should be on a new line):

```
DA2ntYZUKOYF+jweTUt2tY24fw4bk4lZOmrHFNTkdm04MFPM9j3Gqg3h51hA4EnnEo/
DA2ntYZUKOYF+jweTUt2tY24fw4bk4lZOmrHFNTkdm04MFPM9j3Gqg3h51hA4EnnEo/
DA2ntYZUKOYF+jweTUt2tY24fw4bk4lZOmrHFNTkdm04MFPM9j3Gqg3h51hA4EnnEo/
```


The following example illustrates how you can run the CLI tool for bulk update auth tokens and grant access to list of pipes:

```bash
tinybird-cli tokens put \
 -a=p.eyJ1IjogImUwMGU2NjJjLWQ... \
 -f=/path/to/encrypted-tokens.txt \
 -d=base64:FCJXJ1i835dXnGk07Nu1rOLl...
 -s=PIPES:READ:PipeName,PIPES:READ:PipeName2,PIPES:READ:PipeName3 \ 
```


## Build
```
go build
```

## TODO
tests, build artifacts



