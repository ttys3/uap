# ua-proxy

User Agent proxy for windows guest machine

redirect all https/http/ftp protocol handler from your windows guest machine to the host machine

## run `uapd` on Linux host machine

open `uapd.service` change `UAP_AUTH` to your own password

```bash
make install
```

## run `uap` on windows guest machine

open setup.bat change `UAP_AUTH` to the same as `uapd`

run setup.bat
run clients.reg,  registered.reg,  UAPHTML.reg

`Settings` -  `Default Apps` - `Web Browser` - select `UAP Explorer`