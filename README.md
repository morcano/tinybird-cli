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
| --scope      | -s  | This flag is used to provide a list of scopes ([for tokens API](https://www.tinybird.co/docs/api-reference/token-api#put--v0-tokens-(.+))) |

## Usage
According to directory structure app args looks as follows
```
tinybird-cli [entity] [action] [flags]
```


The following example illustrates how you can run the CLI tool for bulk update auth tokens and grant access to list of pipes:

```bash
tinybird-cli tokens put \
 -a=p.eyJ1IjogImUwMGU2NjJjLWQ... \
 -s=PIPES:READ:PipeName,PIPES:READ:PipeName2,PIPES:READ:PipeName3 \ 
```


## Build
```
go build
```

## TODO
tests, build artifacts



