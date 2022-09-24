# GetState message

## Examples
000:  de 14 00 00 00 00 (between SetPreamblePattern and first de16 message)
000:  de 15 00 00 00 00 Idle message
000:  de 16 00 00 00 00 Normal message
000:  de 0b 00 00 00 00 (detected via USB sniffer)

## Fields
00:    messageID
01:    stateID
02-05: 00
