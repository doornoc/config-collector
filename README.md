# Config Collector

## 実行方法
- execute one  
```./backend start  --config [config path] --template [template config path]```
- Execute periodically   
```./backend start  --config [config path] --template [template config path]```

## Template
```
- os_type:      OS Type
- commands:     Execute command
- ignore_line:  When ignoring strings
- config_start: Where the config start
- config_end:   Where the config end
```