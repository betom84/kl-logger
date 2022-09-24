### SetFrame

#### Actions
00: rtGetHistory     - Ask for History message
01: rtSetTime        - Ask for Send Time to weather station message
02: rtSetConfig      - Ask for Send Config to weather station message
02: rtReqFirstConfig - Ask for Send (First) Config to weather station message
03: rtGetConfig      - Ask for Config message
04: rtGetCurrent     - Ask for Current Weather message
20: Send Config      - Send Config to WS
60: Send Time        - Send Time to WS (only if station is just initialized)

000:  d5 00 0b DevID LI 00 CfgCS 8cINT ThisAddr xx xx xx  rtGetHistory
000:  d5 00 0b DevID LI 01 CfgCS 8cINT ThisAddr xx xx xx  rtReqSetTime
000:  d5 00 0b DevID LI 02 CfgCS 8cINT ThisAddr xx xx xx  rtReqSetConfig
000:  d5 00 0b f0 f0 ff 03 ff ff 8cINT DevID LI xx xx xx  rtReqFirstConfig
000:  d5 00 0b DevID LI 03 CfgCS 8cINT ThisAddr xx xx xx  rtGetConfig
000:  d5 00 0b DevID LI 04 CfgCS 8cINT ThisAddr xx xx xx  rtGetCurrent
000:  d5 00 7d DevID LI 20 [ConfigData  .. .. .. .. CfgCS] Send Config
000:  d5 00 0d DevID LI 60 CfgCS [TimeData .. .. .. .. ..  Send Time

##### All SetFrame messages:
00:    messageID
01:    00
02:    Message length (starting with next byte)
03-04: DeviceID           [DevID]
05:    LI/ff              Logger ID: 0-9 = Logger 1 - logger 10
06:    Action
07-08: Config checksum    [CfgCS]

##### Additional bytes rtGetCurrent, rtGetHistory, rtSetTime messages:
09hi:    0x80               (meaning unknown, 0.5 byte)
09lo-10: ComInt             [cINT]    1.5 byte
11-13:   ThisHistoryAddress [ThisAddr] 3 bytes (high byte first)

##### Additional bytes Send Time message:
09:    seconds
10:    minutes
11:    hours
12hi:  day_lo         (low byte)
12lo:  DayOfWeek      (mo=1, tu=2, we=3, th=4, fr=5, sa=6 su=7)
13hi:  month_lo       (low byte)
13lo:  day_hi         (high byte)
14hi:  (year-2000)_lo (low byte)
14lo:  month_hi       (high byte)
15hi:  not used
15lo:  (year-2000)_hi (high byte)

WS SetTime - Send time to WS
Time  d5 00 0d 01 07 00 60 1a b1 25 58 21 04 03 41 01
time sent: Thu 2014-10-30 21:58:25

### SendConfig - ConfigData

Frame offset    : 08 bytes
Message length  : 0x7d (125 bytes)
Message offset  : 05 bytes

| Offset | Nibbles | Description
| :----: | :-----: | ---
| 0      | 2       | Settings 8=? | 0-7=contrast, 8=alert OFF, 4=DCF ON, 2=clock 12h, 1=temp-F
| 1      | 2       | TimeZone difference with Frankfurt (CET) f4 (-12) =tz -12h, 00=tz 0h, 0c (+12) = tz +12h
| 2      | 2       | History Interval
| 3      | 3       | Temp0Max (reverse group 1)
| 4,5    | 3       | Temp0Min (reverse group 1)
| 6      | 3       | Temp1Max (reverse group 2)
| 7,5    | 3       | Temp1Min (reverse group 2)
| 9      | 3       | Temp2Max (reverse group 3)
| 10,5   | 3       | Temp2Min (reverse group 3)
| 12     | 3       | Temp3Max (reverse group 4)
| 13,5   | 3       | Temp3Min (reverse group 4)
| 15     | 3       | Temp4Max (reverse group 5)
| 16,5   | 3       | Temp4Min (reverse group 5)
| 18     | 3       | Temp5Max (reverse group 6)
| 19,5   | 3       | Temp5Min (reverse group 6)
| 21     | 3       | Temp6Max (reverse group 7)
| 22,5   | 3       | Temp6Min (reverse group 7)
| 24     | 3       | Temp7Max (reverse group 8)
| 25,5   | 3       | Temp7Min (reverse group 8)
| 27     | 3       | Temp8Max (reverse group 9)
| 28,5   | 3       | Temp8Min (reverse group 9)
| 30     | 2       | Humidity0Max (reverse group 10)
| 31     | 2       | Humidity0Min (reverse group 10)
| 32     | 2       | Humidity1Max (reverse group 11)
| 33     | 2       | Humidity1Min (reverse group 11)
| 34     | 2       | Humidity2Max (reverse group 12)
| 35     | 2       | Humidity2Min (reverse group 12)
| 36     | 2       | Humidity3Max (reverse group 13)
| 37     | 2       | Humidity3Min (reverse group 13)
| 38     | 2       | Humidity4Max (reverse group 14)
| 39     | 2       | Humidity4Min (reverse group 14)
| 40     | 2       | Humidity5Max (reverse group 15)
| 41     | 2       | Humidity5Min (reverse group 15)
| 42     | 2       | Humidity6Max (reverse group 16)
| 43     | 2       | Humidity6Min (reverse group 16)
| 44     | 2       | Humidity7Max (reverse group 17)
| 45     | 2       | Humidity7Min (reverse group 17)
| 46     | 2       | Humidity8Max (reverse group 18)
| 47     | 2       | Humidity8Min (reverse group 18)
| 48     | 0       | '0000000000' sens0: 8=tmp lo al, 4=tmp hi al, 2=hum lo al, 1=hum hi al; same for sens1-8, 0000
| 53     | 6       | Description1 (reverse)
| 61     | 6       | Description2 (reverse)
| 69     | 6       | Description3 (reverse)
| 77     | 6       | Description4 (reverse)
| 85     | 6       | Description5 (reverse)
| 93     | 6       | Description6 (reverse)
| 101    | 6       | Description7 (reverse)
| 109    | 6       | Description8 (reverse)
| 117    | 2       | '00' (output only) 0000, 1=reset hi-lo values
| 119    | 2       | outBufCS
| 120    | 0       | end

#### History Interval
| Constant | Value | Message received at
| -------- | :---: | ---
| hi01Min  | 0     | 00:00, 00:01, 00:02, 00:03 ... 23:59
| hi05Min  | 1     | 00:00, 00:05, 00:10, 00:15 ... 23:55
| hi10Min  | 2     | 00:00, 00:10, 00:20, 00:30 ... 23:50
| hi15Min  | 3     | 00:00, 00:15, 00:30, 00:45 ... 23:45
| hi30Min  | 4     | 00:00, 00:30, 01:00, 01:30 ... 23:30
| hi01Std  | 5     | 00:00, 01:00, 02:00, 03:00 ... 23:00
| hi02Std  | 6     | 00:00, 02:00, 04:00, 06:00 ... 22:00
| hi03Std  | 7     | 00:00, 03:00, 09:00, 12:00 ... 21:00
| hi06Std  | 8     | 00:00, 06:00, 12:00, 18:00