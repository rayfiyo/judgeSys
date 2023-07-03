# What's this?
The original judging system, which became personally necessary, is now somewhat public.
Not sure if it is useful for anyone but me.
個人的に必要になったオリジナルジャッジシステムをなんとなく公開。
私以外にとって役に立つかは不明。

---
# Usage

# makeAns
以下のディレクトリ構成を含むとき
~~~
ac.c makeAns sample.txt
~~~
./ansMaker を実行すると ans.txt が生成される。
このテキストファイルの中身は、 sample.txt を入力としたときの ac.c の出力結果である。

# compile
以下のディレクトリ構成を含むとき
~~~
ans.txt compile sample.txt 2023-06-12.c
~~~
./compile を実行すると 2023-06-12.txt が生成される。
このテキストファイルの中身は、 sample.txt を入力したときの 2023-06-12.c の出力結果 や エラー 実行結果(0/1) などである。
