# fedit

fedit renames the files📝 in the specified directory📁. 
And you can edit with favorite editor, for example vim, emacs, notepad and gedit.

If you takes time with rename tasks, recommend this.

<p align="center">
    <img src="./image/fedit-sample.gif" width="200px">
</p>

## Install

> **NOTE** You must have already installed Go.

```
go get github.com/yasukotelin/fedit
```

## How to use

```
> fedit ./
```

When you run `fedit` with specifying the directory, the file list is opened by **Default editor**.
Default editor is notepad when uses on the windows, and it's vim when on the other os.

### Specify editor

You can specify the favorite editor with `-e` option.

```
fedit -e gvim ./
```


### Help

```
> fedit --help
```
