polling_intervall = 10s
comm_intervall = 8

first_sleep = 0.3s (300ms)

firstSleep = 1
nextSleep = 1

doRF
- while running doRFCommunication

doRFSetup
- firstSleep = 0.075
- nextSleep = 0.005

doRFCommunication
- sleep(firstSleep)
- getState err -> 5s sleep
- state != 0x16 -> sleep(nextSleep) -> getState
- getFrame
- generateResponse
    - DataWritten -> setRX, raise DataWritten
    - GetConfig -> nextSleep = 0.010
    - GetCurrent -> nextSleep = 0.010
        - e.g. get/set config
        - ask for History instead of Current (kl.py#3633)
    - GetHistory -> nextSleep = 0.010
        - ask for GetHistory with nextIndex
    - ResponseRequest
        - ReqReadHistory -> setSleep(0.075, 0.005)
        - ReqFirstConfig -> firstResponseFrame + setSleep(0.075, 0.005) 
        - ReqSetConfig -> configFrame + setSleep(0.075, 0.005)
        - ReqSetTime -> timeFrame + setSleep(0.075, 0.005)
        - default -> getHistory + nextSleep = 0.010
    - default -> raise BadResponse
    - unknown device -> setSleep(0.200, 0.005)
- setFrame + setTx, err (DataWritten, BadResponse) -> setRx

```plantuml
activate Application
...sleep...

Application->Transceiver: getState
return state

Application->Transceiver: getFrame
return requestFrame

Application->Application++: processFrame
return responseFrame

Application->Transceiver: setFrame
Application->Transceiver: setTx
Application->Transceiver: setRx
deactivate Application
...sleep...
```

@import "flow.puml"