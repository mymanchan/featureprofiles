# gNOI-5.3: Copying Debug Files

## Summary

Validate that the debug files can be copied out of the DUT.

## Procedure

* Issue KillProcessRequest to the DUT to crash a software process.
* Issue gnoi.Healthz Get RPC to chassis.
* Validate that the device returns the vendor relevant information for
    debugging.

## Config Parameter Coverage

N/A

## Telemetry Parameter Coverage

N/A

## Protocol/RPC Parameter Coverage

* gNOI
  * Healthz
  * System
    * KillProcessRequest
