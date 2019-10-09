# guicast

![guicast](https://github.com/yasutakatou/guicast/blob/pic/guicast.gif)

WindowsのGUIをプロンプトから連続でキャプチャとりつつ自動的に操作するCLIツール

・以前作ったのはキャプチャとキーエミュレートが外部コマンドに依存していたのでWin32API使ってネイティブで動作するようにした
・robotgoはzlib1.dllが必要でgolangの良さであるワンバイナリが崩れるので根性で他のライブラリ組み込んだ
　（あとrobotgoはLinuxからWindowsみたいなクロスコンパイルがどうしてもできないし）
・Win32APIでWindowタイトルの一覧作って、キーワードにあてはまるのだけHwndを持ってきてフォアグランドに出す実装が無かったので。。
https://stackoverflow.com/questions/47189825/golang-how-to-set-window-on-top
→「robotgoは正確にactive window設定できないぜ」というこのやりとりしかない。ので根性実装

動作として指定したwindowタイトルに同じキー操作を投げ込んでキャプチャを取るの繰り返し
