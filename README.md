# Config Collector

## 実行方法
- 一回だけ  
```./backend start once --config [config path] --template [template config path]```
- 定期的に起動   
```./backend start cron --config [config path] --template [template config path]```
- Execute on container  
```cd docker && mkdir config && cp ../config/* config/ && docker compose up -d```
### 動作確認(github commit & push機能を無効)
- 一回だけ  
```./backend test once --config [config path] --template [template config path]```
- 定期的に起動   
```./backend test cron --config [config path] --template [template config path]```


## Template
sample templateはconfig/template.jsonにあります
```
- os_type:       OS Type
- commands:      実行するコマンド(配列順に実行)
- config_start:  config判定開始時の文字列(一部一致した次の行から抽出を開始)
- config_end:    config判定終了時の文字列(一部一致した次の行から抽出を終了)
- ignore_line:   必要のないconfig行を削除(一部一致した行がconfigから削除)
- input_console: 特殊オプション
```
### その他
Templateのcommands項目にて ```{{```と```}}```を使用することで、任意の値(パスワードなど)に置き換えることが可能です。  
Ciscoなどの```enable password```などに対応することが可能となっております。  
また、現時点では置き換え(replace)にしか対応していませんが、 将来的には、if文などへの対応も検討しております。

#### 例) 
Templateに```{{ enable_password }}```のように利用すると、config.json内のoptionから  
```
option: {
    enable_password: "password_here"
}
```
enable_passwordのキーに対して必要な値を返す仕組みになっています。     


## Config
sample templateはconfig/config.jsonにあります
### Global
```
- timezone:      Go言語のTimeZone
- github:        ConfigのPush先Githubの設定
- slack_webhook: エラー時に飛ばすSlackのWebhookURL
- tmp_path:      一時ファイル置き場(gitや内部コンフィグの一時置き場として使用)
- exec_time:     定期実行時の周期(秒)
- debug:         Debugモード
```
### Option(自由記述)
```
- enable_password   enableパスワードの記入欄
...
```
### Device
```
- name:       ルータのホスト名(gitにアップロードする際のファイル名になる)
- hostname:   各種機器のSSH時のIPアドレス又はホスト名
- port:       各種機器のSSH時のポート
- user:       各種機器のSSH時のユーザ
- password:   各種機器のSSH時のパスワード
- os_type:    templateに対応するos type
```
