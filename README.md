# feditelin
File edit telin - ☆

ディレクトリ内のファイルを一気にリネームできるツールです。
お好みのエディタでファイル一覧を開いて編集できます。

ファイラーでぽちぽちリネームするのが面倒なときにどうぞ。

![preview](./image/feditelin-prev1.gif)

## Install

```
go get github.com/yasukotelin/feditelin
```

## How to use

```
$ feditelin --help
feditelin is the tool to rename all files in directory

Usage:
  feditelin [flags]

Flags:
  -e, --editor string   specify the editor to open.  (default "notepad")
  -h, --help            help for feditelin
      --version         version for feditelin


```

引数にディレクトリを指定することで、そのディレクトリ直下のファイル一覧情報がエディタで開かれます。そのままリネームして保存、閉じるとリネーム実行の再確認が表示されるので、y or nで確定してください。

### エディタ指定

`-e` オプションで展開するエディタを指定できます。

```
feditelin -e gvim ./
```

また、デフォルトのエディタは、Windowsの場合メモ帳、それ以外の場合はvimに指定されています。

## Develop dependencies

* Go
* dep