# docker-nutcheck

Docker container including tools for Nut UPS

tools included :

- nutcheck : go monitoring daemon that call nut and send telegram notification on battery level change (50/20 %)

## Usage :

### nutcheck
`docker run --rm --privileged -it akit042/docker-nutcheck nutcheck -telegramtoken XXXXXXX:XXXXXXXXXXXXXXX -telegramid YYYYYYYY -nutaddress 192.168.1.252 -nutuser user -nutpassword "secret" -nutupsname ups`

```
Usage of nutcheck:
  -d	debug mode
  -l	List Nut Ups
  -nutaddress string
    	Nut server address (or use env variable : NUT_ADDRESS) (default "127.0.0.1")
  -nutpassword string
    	Nut server password (or use env variable : NUT_PASSWORD)
  -nutport int
    	Nut server port (or use env variable : NUT_PORT) (default 3493)
  -nutupsname string
    	Nut ups to check (or use env variable : NUT_UPS_NAME)
  -nutuser string
    	Nut server user (or use env variable : NUT_USER)
  -poolinginterval int
    	Pooling Interval (or use env variable : POOLING_INTERVAL) (default 30)
  -telegramid int
    	To find an id, please contact @myidbot on telegram (or use env variable : TELEGRAM_ID)
  -telegramtoken string
    	To create a bot, please contact @BotFather on telegram (or use env variable : TELEGRAM_TOKEN) (default "None")
  -v	Print build id
```
