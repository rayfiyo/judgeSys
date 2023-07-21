# What's this?
The original judging system, which became personally necessary, is now somewhat public.<br>
Not sure if it is useful for anyone but me.<br>
個人的に必要になったオリジナルジャッジシステムをなんとなく公開。<br>
私以外にとって役に立つかは不明。<br>

---
# Usage

## generator
以下のディレクトリ構成を含むとき
~~~
ac.c generator sample.txt
~~~
./makeAns を実行すると ans.txt が生成される。
このテキストファイルの中身は、 sample.txt を入力としたときの ac.c の出力結果である。

## compile
以下のディレクトリ構成を含むとき
~~~
ans.txt judge sample.txt 2023-06-12.c
~~~
./compile を実行すると 2023-06-12.txt が生成される。
このテキストファイルの中身は、 sample.txt を入力したときの 2023-06-12.c の出力結果 や エラー 実行結果(0/1) などである。
ans.txt と比較し、実行結果も記載する
