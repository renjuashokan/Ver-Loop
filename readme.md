# Adios...

## Building this project

```
cd <to this project directory>
go build
```

## About Env Variables

**VERLOOP_DEBUG** \
\
Based on this variable log level is set. If it is empty "INFO" is set by defalt.\
\
\
**VERLOOP_DSN** \
I'm assuming you are running postgres docker present in other directory. (for more info, check the readme file on the postgres_renju directory)
Only postgres server ip is needed here. If you want to connect to any other pg db, please edit config.yaml file accordingly
