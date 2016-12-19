# vcprompt

vcprompt is a simple Go program that prints version control system informations. It is designed to
be used by shell prompts.

vcprompt is originally a C program, which happens to be reside in
[here](http://hg.gerg.ca/vcprompt/).

## Installation

```sh
go get -u github.com/igungor/vcprompt
```

## Usage

My Zsh prompt:

```sh
PROMPT='%F{blue}%1~%F{242} %F{yellow}$(vcprompt -f "%n:%b:%m") %F$'
```
