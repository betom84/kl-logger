# ReadConfigFlash
 
## In - receive data

### Example
0000: dc 0a 01 f5 00 01 8d 18 01 02 12 01 0d 01 07 ff ff ff ff ff 00 - freq correction
0000: dc 0a 01 f9 01 02 12 01 0d 01 07 ff ff ff ff ff ff ff ff ff 00 - transceiver data

### Fields
00:    messageID
01:    length
02-03: address

#### Additional bytes frequency correction
05lo-07hi: frequency correction

#### Additional bytes transceiver data
05-10:     serial number
09-10:     DeviceID [devID]

## Out - ask for data

### Example
000: dd 0a 01 f5 58 d8 34 00 90 10 07 01 08 f2 ee - Ask for freq correction
000: dd 0a 01 f9 cc cc cc cc 56 8d b8 00 5c f2 ee - Ask for transceiver data

### Fields
00:    messageID
01:    length
02-03: address
04-14: cc


ReadConfigFlash data

Ask for frequency correction
rcfo  0000: dd 0a 01 f5 cc cc cc cc cc cc cc cc cc cc cc
      0000: dd 0a 01 f5 58 d8 34 00 90 10 07 01 08 f2 ee - Ask for freq correction

readConfigFlash frequency correction
rcfi  0000: dc 0a 01 f5 00 01 78 a0 01 02 0a 0c 0c 01 2e ff ff ff ff ff
      0000: dc 0a 01 f5 00 01 8d 18 01 02 12 01 0d 01 07 ff ff ff ff ff 00 - freq correction
frequency correction: 96416 (0x178a0)
adjusted frequency: 910574957 (3646456d)

Ask for transceiver data
rcfo  0000: dd 0a 01 f9 cc cc cc cc cc cc cc cc cc cc cc
      0000: dd 0a 01 f9 cc cc cc cc 56 8d b8 00 5c f2 ee - Ask for transceiver data

readConfigFlash serial number and DevID
rcfi  0000: dc 0a 01 f9 01 02 0a 0c 0c 01 2e ff ff ff ff ff ff ff ff ff
      0000: dc 0a 01 f9 01 02 12 01 0d 01 07 ff ff ff ff ff ff ff ff ff 00 - transceiver data
transceiver ID: 302 (0x012e)
transceiver serial: 01021012120146

