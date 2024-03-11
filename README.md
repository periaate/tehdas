# Tehdas
Tehdas is an overengineered build automation tool developed to make building my programs a bit more convenient.

## Usage
Running `tehdas` will create both the files and the directories it uses to work. This includes `~/.tehdas/` directory and the `~/.tehdas/.conf` file. If the `.conf` file doesn't specify a target, the program will prompt you for the desired target path which will be used as default after this configuration process. Inputting nothing will assume the default target, which is `~/.tehdas/`.

Having configured `tehdas`, the program will infer the project language and build the project accordingly, placing the binary to the configured target folder.

Calling `tehdas` syntax is as follows:\
`tehdas {entry path}`\
Entry path may be a file or directory. By default the entry is the current working directory.

## Why another build tool?
Sorry, I procrastinate by writing software.

## Supported Languages
- Go
- I'll probably add support for Rust Someday™.

## Additional notes
This project is overengineered, it isn't idiomatic, nor does it try to be. I made this for fun, so I might have used patterns which aren't recommended or are considered bad practices. It implements a custom key-value string coding format for its configuration file (which only ever has a single key and a single value). It uses nonstandard error handling, control flow, and logging.

## License
This project is licensed under [EUPL 1.2.](./LICENSE)