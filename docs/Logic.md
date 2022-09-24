Program Logic

The RF communication thread uses the following logic to communicate with the
weather station console:

Step 1.  Perform in a while loop getState commands until state 0xde16
         is received.

Step 2.  Perform a getFrame command to read the message data.

Step 3.  Handle the contents of the message. The type of message depends on
         the response type:

  Response type (hex):
  10: WS SetTime / SetConfig - Data written
      confirmation the setTime/setConfig setFrame message has been received
      by the console
  20: GetConfig
      save the contents of the configuration for later use (i.e. a setConfig
      message with one ore more parameters changed)
  30: Current Weather
      handle the weather data of the current weather message
  40: Actual / Outstanding History
      ignore the data of the actual history record when there is no data gap;
      handle the data of a (one) requested history record (note: in step 4 we
      can decide to request another history record).
  50: Request Read-History (MEM % > 0)
      no other action than debug log (the driver will always read history
      messages when available)
  51: Request First-Time Config
      prepare a setFrame first time message
  52: Request SetConfig
      prepare a setFrame setConfig message
  53: Request SetTime
      prepare a setFrame setTime message

Step 4.  When  you  didn't receive the message in step 3 you asked for (see
         step 5 how to request a certain type of message), decide if you want
         to ignore or handle the received message. Then go to step 5 to
         request for a certain type of message unless the received message
         has response type 51, 52 or 53, then prepare first the setFrame
         message the wireless console asked for.

Step 5.  Decide what kind of message you want to receive next time. The
         request is done via a setFrame message (see step 6).  It is
         not guaranteed that you will receive that kind of message the next
         time but setting the proper timing parameters of firstSleep and
         nextSleep increase the chance you will get the requested type of
         message.

Step 6. The action parameter in the setFrame message sets the type of the
        next to receive message.

  Action (hex):

  00: rtGetHistory - Ask for History message
                     setSleep(FIRST_SLEEP, 0.010)
  01: rtSetTime    - Ask for Send Time to weather station message
                     setSleep(0.075, 0.005)
  02: rtSetConfig  - Ask for Send Config to weather station message
                     setSleep(FIRST_SLEEP, 0.010)
  03: rtGetConfig  - Ask for Config message
                     setSleep(0.400, 0.400)
  04: rtGetCurrent - Ask for Current Weather message
                     setSleep(FIRST_SLEEP, 0.010)
  20: Send Config  - Send Config to WS
                     setSleep(0.075, 0.005)
  60: Send Time    - Send Time to WS
                     setSleep(0.075, 0.005)

  Note: after the Request First-Time Config message (response type = 0x51)
        perform a rtGetConfig with setSleep(0.075,0.005)

Step 7. Perform a setTX command

Step 8. Go to step 1 to wait for state 0xde16 again.

## Messages

ID | Name                | Length (bytes)
-- | ------------------- | --------------
00 | GetFrame            | 0x111 (273)
d0 | SetRX               | 0x15  (21)
d1 | SetTX               | 0x15  (21)
d5 | SetFrame            | 0x111 (273)
d7 | SetState            | 0x15  (21)
d8 | SetPreamblePattern  | 0x15  (21)
d9 | Execute             | 0x0f  (15)
dc | ReadConfigFlash<    | 0x15  (21)
dd | ReadConfigFlash>    | 0x15  (21)
de | GetState            | 0x0a  (10)
f0 | WriteReg            | 0x05  (5)