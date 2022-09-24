# GetFrame

## Fields
| Byte   | Len | Symbol   | Description          
| :----: | :-: | :------: | ---
| 00     | 01  | 00       | GetFrame ID
| 01     | 01  | 00       | const                    
| 02     | 01  | hex      | Message Length (starting with next byte)
| 03     | 02  | DevID    | Device ID / `f0 f0` at init
| 05     | 01  | LI       | Logger ID 0-9 / `ff` at init
| 06     | 01  | hex      | Message ID
| 07     | 01  | SQ       | Signal Quality (in steps of 5) 
| 07[^1] | 01  | MP       | Memory Percentage not read to server (in steps of 5)
| 08     | 02  | CfgCS    | Config checksum
| 10[^2] | 03  | LateAddr | Latest History Record = (LateAddr - `0x070000`) / 32 
| 13[^2] | 03  | ThisAddr | This History Record = (ThisAddr - `0x070000`) / 32 

## Frames
| Type                      | ID     | Format                                                        |
| ------------------------- | :----: | ------------------------------------------------------------- |
| DataWritten               | `0x10` | `00 00 07 DevID LI 10 SQ CfgCS xx xx xx xx xx xx xx xx xx xx` |
| GetConfig                 | `0x20` | `00 00 7d DevID LI 20 SQ [ConData .. .. .. .. .. .. .] CfgCS` |
| Current weather data      | `0x30` | `00 00 e5 DevID LI 30 SQ CfgCS [CurData .. .. .. .. .. .. .]` |
| History weather data      | `0x40` | `00 00 b5 DevID LI 40 SQ CfgCS LateAddr ThisAddr [HisData .]` |
| Request to read history   | `0x50` | `00 00 b5 DevID LI 50 MP CfgCS xx xx xx xx xx xx xx xx xx xx` |
| Request initial config    | `0x51` | `00 00 07 f0 f0 ff 51 SQ CfgCS xx xx xx xx xx xx xx xx xx xx` |
| Request SetConfig         | `0x52` | `00 00 07 DevID LI 52 SQ CfgCS xx xx xx xx xx xx xx xx xx xx` |
| Request SetTime           | `0x53` | `00 00 07 DevID LI 53 SQ CfgCS xx xx xx xx xx xx xx xx xx xx` |

### GetConfig (ConData)

Frame offset        : 08 bytes
Message length      : 0x7d (125 bytes)
Message data offset : 05 bytes

| Offset | Nibbles | Description
| :----: | :-----: | ---
| 0      | 2       | Settings
| 1      | 2       | TimeZone
| 2      | 2       | HistoryInterval
| 3      | 3       | Temperature KlimaLogg Max
| 4,5    | 3       | Temperature KlimaLogg Min
| 6      | 3       | Temperature Sensor 1 Max
| 7,5    | 3       | Temperature Sensor 1 Min
| 9      | 3       | Temperature Sensor 2 Max
| 10,5   | 3       | Temperature Sensor 2 Min
| 12     | 3       | Temperature Sensor 3 Max
| 13,5   | 3       | Temperature Sensor 3 Min
| 15     | 3       | Temperature Sensor 4 Max
| 16,5   | 3       | Temperature Sensor 4 Min
| 18     | 3       | Temperature Sensor 5 Max
| 19,5   | 3       | Temperature Sensor 5 Min
| 21     | 3       | Temperature Sensor 6 Max
| 22,5   | 3       | Temperature Sensor 6 Min
| 24     | 3       | Temperature Sensor 7 Max
| 25,5   | 3       | Temperature Sensor 7 Min
| 27     | 3       | Temperature Sensor 8 Max
| 28,5   | 3       | Temperature Sensor 8 Min
| 30     | 2       | Humidity KlimaLogg Max
| 31     | 2       | Humidity KlimaLogg Min
| 32     | 2       | Humidity Sensor 1 Max
| 33     | 2       | Humidity Sensor 1 Min
| 34     | 2       | Humidity Sensor 2 Max
| 35     | 2       | Humidity Sensor 2 Min
| 36     | 2       | Humidity Sensor 3 Max
| 37     | 2       | Humidity Sensor 3 Min
| 38     | 2       | Humidity Sensor 4 Max
| 39     | 2       | Humidity Sensor 4 Min
| 40     | 2       | Humidity Sensor 5 Max
| 41     | 2       | Humidity Sensor 5 Min
| 42     | 2       | Humidity Sensor 6 Max
| 43     | 2       | Humidity Sensor 6 Min
| 44     | 2       | Humidity Sensor 7 Max
| 45     | 2       | Humidity Sensor 7 Min
| 46     | 2       | Humidity Sensor 8 Max
| 47     | 2       | Humidity Sensor 8 Min
| 48     | 10      | AlarmSet
| 53     | 16      | Description Sensor 1
| 61     | 16      | Description Sensor 2
| 69     | 16      | Description Sensor 3
| 77     | 16      | Description Sensor 4
| 85     | 16      | Description Sensor 5
| 93     | 16      | Description Sensor 6
| 101    | 16      | Description Sensor 7
| 109    | 16      | Description Sensor 8
| 117    | 2       | ResetHiLo (output only)
| 119    | 2       | inBufCS
| 120    | 0       | end

#### AlarmSet

| Alarm                     | Value
| :------------------------ | ---
| Humidity KlimaLogg Max    | `00 00 00 00 01`
| Humidity KlimaLogg Min    | `00 00 00 00 02`
| Temperature KlimaLogg Max | `00 00 00 00 04`
| Temperature KlimaLogg Min | `00 00 00 00 08`
| Humidity Sensor 1 Max     | `00 00 00 00 10`
| Humidity Sensor 1 Min     | `00 00 00 00 20`
| Temperature Sensor 1 Max  | `00 00 00 00 40`
| Temperature Sensor 1 Min  | `00 00 00 00 80`
| ...                       | ...

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

### Current Weather (CurData)

Frame offset        : 10 bytes
Message length      : 0xe5 (229 bytes)
Message data offset : 07 bytes

| Offset | Nibbles | Description
| :----: | :-----: | ---
| 0      | 58      | Weather KlimaLogg
| 24     | 58      | Weather Sensor 1
| 48     | 58      | Weather Sensor 2
| 72     | 58      | Weather Sensor 3
| 96     | 58      | Weather Sensor 4
| 120    | 58      | Weather Sensor 5
| 144    | 58      | Weather Sensor 6
| 168    | 58      | Weather Sensor 7
| 192    | 58      | Weather Sensor 8
| 216    | 12      | AlarmData (`00 00 00 00 00 00`)
| 222    | 0       | end

#### Weather record

| Offset | Nibbles | Description
| :----: | :-----: | ---
| 0      | 8       | Humidity Max DateTime
| 4      | 8       | Humidity Min DateTime
| 8      | 2       | Humidity Max
| 9      | 2       | Humidity Min
| 10     | 2       | Humidity
| 11     | 1       | '0'
| 11.5   | 8       | Temperature Max DateTime
| 15.5   | 8       | Temperature Min DateTime
| 19.5   | 3       | Temperature Max
| 21     | 3       | Temperature Min
| 22.5   | 3       | Temperature

#### AlarmData 

Length: 6 bytes (12 nibbles of 4 bit)
Nibbles: 01 02 03 04 05 06 07 08 09 10 11 12 (meaning of 01-06 unknown)

| 07 | 08 | 09 | 10 | 11 | 12 | Value
| -- | -- | -- | -- | -- | -- | ---
| 80 | 00 | 00 | 00 | 00 | 00 | Sensor 8 TX batt low
| 40 | 00 | 00 | 00 | 00 | 00 | Sensor 7 TX batt low
| 20 | 00 | 00 | 00 | 00 | 00 | Sensor 6 TX batt low
| 10 | 00 | 00 | 00 | 00 | 00 | Sensor 5 TX batt low
| 08 | 00 | 00 | 00 | 00 | 00 | Sensor 4 TX batt low
| 04 | 00 | 00 | 00 | 00 | 00 | Sensor 3 TX batt low
| 02 | 00 | 00 | 00 | 00 | 00 | Sensor 2 TX batt low
| 01 | 00 | 00 | 00 | 00 | 00 | Sensor 1 TX batt low
| 00 | 80 | 00 | 00 | 00 | 00 | KlimaLogg RX batt low
| 00 | 40 | 00 | 00 | 00 | 00 |
| 00 | 20 | 00 | 00 | 00 | 00 |
| 00 | 10 | 00 | 00 | 00 | 00 |
| 00 | 08 | 00 | 00 | 00 | 00 | Temperature Sensor 8 Min
| 00 | 04 | 00 | 00 | 00 | 00 | Temperature Sensor 8 Max
| 00 | 02 | 00 | 00 | 00 | 00 | Humidity Sensor 8 Min
| 00 | 01 | 00 | 00 | 00 | 00 | Humidity Sensor 8 Max
| 00 | 00 | 80 | 00 | 00 | 00 | Temperature Sensor 7 Min
| 00 | 00 | 40 | 00 | 00 | 00 | Temperature Sensor 7 Max
| 00 | 00 | 20 | 00 | 00 | 00 | Humidity Sensor 7 Min
| 00 | 00 | 10 | 00 | 00 | 00 | Humidity Sensor 7 Max
| 00 | 00 | 08 | 00 | 00 | 00 | Temperature Sensor 6 Min
| 00 | 00 | 04 | 00 | 00 | 00 | Temperature Sensor 6 Max
| 00 | 00 | 02 | 00 | 00 | 00 | Humidity Sensor 6 Min
| 00 | 00 | 01 | 00 | 00 | 00 | Humidity Sensor 6 Max
| 00 | 00 | 00 | 80 | 00 | 00 | Temperature Sensor 5 Min
| 00 | 00 | 00 | 40 | 00 | 00 | Temperature Sensor 5 Max
| 00 | 00 | 00 | 20 | 00 | 00 | Humidity Sensor 5 Min
| 00 | 00 | 00 | 10 | 00 | 00 | Humidity Sensor 5 Max 
| 00 | 00 | 00 | 08 | 00 | 00 | Temperature Sensor 4 Min
| 00 | 00 | 00 | 04 | 00 | 00 | Temperature Sensor 4 Max
| 00 | 00 | 00 | 02 | 00 | 00 | Humidity Sensor 4 Min
| 00 | 00 | 00 | 01 | 00 | 00 | Humidity Sensor 4 Max
| 00 | 00 | 00 | 00 | 80 | 00 | Temperature Sensor 3 Min
| 00 | 00 | 00 | 00 | 40 | 00 | Temperature Sensor 3 Max
| 00 | 00 | 00 | 00 | 20 | 00 | Humidity Sensor 3 Min
| 00 | 00 | 00 | 00 | 10 | 00 | Humidity Sensor 3 Max
| 00 | 00 | 00 | 00 | 08 | 00 | Temperature Sensor 2 Min
| 00 | 00 | 00 | 00 | 04 | 00 | Temperature Sensor 2 Max
| 00 | 00 | 00 | 00 | 02 | 00 | Humidity Sensor 2 Min
| 00 | 00 | 00 | 00 | 01 | 00 | Humidity Sensor 2 Max
| 00 | 00 | 00 | 00 | 00 | 80 | Temperature Sensor 1 Min 
| 00 | 00 | 00 | 00 | 00 | 40 | Temperature Sensor 1 Max
| 00 | 00 | 00 | 00 | 00 | 20 | Humidity Sensor 1 Min
| 00 | 00 | 00 | 00 | 00 | 10 | Humidity Sensor 1 Max
| 00 | 00 | 00 | 00 | 00 | 08 | Temperature KlimaLogg Min
| 00 | 00 | 00 | 00 | 00 | 04 | Temperature KlimaLogg Max
| 00 | 00 | 00 | 00 | 00 | 02 | Humidity KlimaLogg Min
| 00 | 00 | 00 | 00 | 00 | 01 | Humidity Klimalogg Max

### History Weather (HisData)

Frame offset        : 16 bytes
Message length      : 0xb5 (181 bytes)
Message data offset : 13 bytes

| Offset | Nibbles | Description
| :----: | :-----: | ---
| 0      | 56      | Record 6 (newest)
| 28     | 56      | Record 5
| 56     | 56      | Record 4
| 84     | 56      | Record 3
| 112    | 56      | Record 2
| 140    | 56      | Record 1 (oldest)
| 168    | 0       | End message

#### Alarm record

| Offset | Nibbles | Description
| :----: | :-----: | ---
| 0      | 26      | Unused
| 13     | 2       | Humidity High
| 14     | 2       | Humidity Low
| 15     | 2       | Humidity
| 16     | 3       | Temperature High
| 17,5   | 3       | Temperature Low
| 19     | 1       | '0'
| 19,5   | 3       | Temperature
| 21     | 1       | Alarmdata; 1=Hum Hi Al, 2=Hum Lo Al, 4=Tmp Hi Al, 8=Tmp Lo Al
| 21,5   | 1       | Sensor; 0=KlimaLogg, 1-8
| 22     | 10      | DateTime
| 27     | 2       | 'ee'

#### History record

| Offset | Nibbles | Description
| :----: | :-----: | ---
| 0      | 2       | Humidity Sensor 8
| 1      | 2       | Humidity Sensor 7
| 2      | 2       | Humidity Sensor 6
| 3      | 2       | Humidity Sensor 5
| 4      | 2       | Humidity Sensor 4
| 5      | 2       | Humidity Sensor 3
| 6      | 2       | Humidity Sensor 2
| 7      | 2       | Humidity Sensor 1
| 8      | 2       | Humidity KlimaLogg
| 9      | 1       | Unused
| 9.5    | 3       | Temperature Sensor 8
| 11     | 3       | Temperature Sensor 7
| 12.5   | 3       | Temperature Sensor 6
| 14     | 3       | Temperature Sensor 5
| 15.5   | 3       | Temperature Sensor 4
| 17     | 3       | Temperature Sensor 3
| 18.5   | 3       | Temperature Sensor 2
| 20     | 3       | Temperature Sensor 1
| 21.5   | 3       | Temperature KlimaLogg
| 23     |10       | DateTime

### Conversion

#### Date (5 bytes)
###### Example (2013-06-21)

| Byte  | Value | Type | Conversion
| ----- | ----- | ---- | ----
| byte1 | 1     | dec  | `year  =  2000 + 10 * byte1`
| byte2 | 3     | dec  | `year  += byte2`
| byte3 | 6     | hex  | `month =  byte3` 
| byte4 | 2     | dec  | `day   =  10 * byte4`
| byte5 | 1     | dec  | `day   += byte5`

#### Time (3 bytes)
###### Example value 00:52

| Byte  | Value | Type | Conversion
| ----- | ----- | ---- | ----
| byte1 | 0     | hex  | `byte1 >= 10 ? hours = 10 + byte1 : hours = byte1`
| byte2 | 5     | hex  | `byte2 >= 10 ? hours += 10; minutes = (byte2 - 10) *10 : minutes = byte2 * 10`
| byte3 | 2     | dec  | `minutes += byte3`

#### DateTime (10 bytes)
###### Example value 2013-05-16 19:15

| Byte   | Value | Conversion
| ------ | ----- | ----
| byte1  | 1     | `year    =  2000 + (byte1 * 10)`
| byte2  | 3     | `year    += byte2`
| byte3  | 0     | `month   =  byte3 * 10`
| byte4  | 5     | `month   += byte4`
| byte5  | 1     | `day     =  byte5 * 10`
| byte6  | 6     | `day     += byte6`
| byte7  | 1     | `hours   =  byte7 * 10`
| byte8  | 9     | `hours   += byte8`
| byte9  | 1     | `minutes =  byte9 * 10`
| byte10 | 5     | `minutes += byte10`

#### Humidity
###### Example value 50

| Byte  | Value | Conversion
| ----- | ----- | ----
| byte1 | 5     | `humidity =  byte1 * 10`
| byte2 | 0     | `humidity += byte2`

#### Temperature conversion
###### Example value 23.2

| Byte  | Value | Conversion
| ----- | ----- | ----
| byte1 | 6     | `temp =  (byte1 * 10) - 40`
| byte2 | 3     | `temp += byte2`
| byte3 | 2     | `temp += (byte3 * 0.1)`


[^1]: Message ID 0x50 (Request to read History) only
[^2]: Message ID 0x40 (History weather data) only