// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

/*
Package litering allows users of gopacket to read packets off the wire or from
PCAP files.

This package is meant to be used with its parent,
http://github.com/google/gopacket, although it can also be used independently
if you just want to get packet data from the wire.

Reading litering Files

The following code can be used to read in data from a litering file.

 if handle, err := litering.OpenOffline("/path/to/my/file"); err != nil {
   panic(err)
 } else {
   packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
   for packet := range packetSource.Packets() {
     handlePacket(packet)  // Do something with a packet here.
   }
 }

Reading Live Packets

The following code can be used to read in data from a live device, in this case
"eth0".

 if handle, err := litering.OpenLive("eth0", 1600, true, 0); err != nil {
   panic(err)
 } else if err := handle.SetBPFFilter("tcp and port 80"); err != nil {  // optional
   panic(err)
 } else {
   packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
   for packet := range packetSource.Packets() {
     handlePacket(packet)  // Do something with a packet here.
   }
 }

Inactive Handles

Newer litering functionality requires the concept of an 'inactive' litering handle.
Instead of constantly adding new arguments to litering_open_live, users now call
litering_create to create a handle, set it up with a bunch of optional function
calls, then call litering_activate to activate it.  This library mirrors that
mechanism, for those that want to expose/use these new features:

  inactive, err := litering.NewInactiveHandle(deviceName)
  if err != nil {
    log.Fatal(err)
  }
  defer inactive.CleanUp()

  // Call various functions on inactive to set it up the way you'd like:
  if err = inactive.SetTimeout(time.Minute); err != nil {
    log.Fatal(err)
  } else if err = inactive.SetTimestampSource("foo"); err != nil {
    log.Fatal(err)
  }

  // Finally, create the actual handle by calling Activate:
  handle, err := inactive.Activate()  // after this, inactive is no longer valid
  if err != nil {
    log.Fatal(err)
  }
  defer handle.Close()

  // Now use your handle as you see fit.

litering Timeouts

litering.OpenLive and litering.SetTimeout both take timeouts.
If you don't care about timeouts, just pass in BlockForever,
which should do what you expect with minimal fuss.

A timeout of 0 is not recommended.  Some platforms, like Macs
(http://www.manpages.info/macosx/litering.3.html) say:
  The read timeout is used to arrange that the read not necessarily return
  immediately when a packet is seen, but that it wait for some amount of time
  to allow more packets to arrive and to read multiple packets from the OS
  kernel in one operation.
This means that if you only capture one packet, the kernel might decide to wait
'timeout' for more packets to batch with it before returning.  A timeout of
0, then, means 'wait forever for more packets', which is... not good.

To get around this, we've introduced the following behavior:  if a negative
timeout is passed in, we set the positive timeout in the handle, then loop
internally in ReadPacketData/ZeroCopyReadPacketData when we see timeout
errors.

litering File Writing

This package does not implement litering file writing.  However, gopacket/literinggo
does!  Look there if you'd like to write litering files.
*/
package litering
