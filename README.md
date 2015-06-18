# vcprompt

vcprompt is a simple Go program that prints version control system informations. It is designed to
be used by shell prompts.

vcprompt is originally a C program, which happens to be reside in
[here](http://hg.gerg.ca/vcprompt/). Original program supports multiple version control systems
-including CVS, SVN, mercurial and git-.

## Installation

```sh
go get github.com/igungor/vcprompt
```

## Usage

My Zsh prompt:

```sh
PROMPT='%F{blue}%1~%F{242} %F{yellow}$(vcprompt -f "%n:%b:%m") %F$'
```
