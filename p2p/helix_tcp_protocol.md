# TCP Protocol for Helix Network

DEFINITIONS:

- User of the Network : Someone who uses the network for its features, this "person" does not run a node and does not deal with any of the consensus and distribution of the network. Their only concern is using the network
- Node of the Network : A "person" whom is involved with the consensus and distribution process.

Ok so the first byte of the message will be used for
checking what type of client are we dealing with,

0x01 = The client sending the TCP message is a user of the network
0x02 = The client sending the TCP message is a node of the network

If 0x01 then the message will be capped at 1024 bytes

0x27899ffaCe558bdE9F284Ba5C8c91ec79EE60FD6

Say we want to do a "Function Call" on the network,
in this case our function call is to move token A from
user a to user b

(Since the Helix Network has a built in token we don't need to use
any external programs)

Biggest Hexadecimal Supported for Amounts 0xFFFFFFFFFF

The TCP should contain the following:

A simplified instruction ->
Payment From 0xC0F2329f646fd2F310862d7dB46451712201EE2e to ->
0x605302D6a2444b0dC08bf107e5FeD12f66bd0059 for 18446744073709551615 (0xFFFFFFFFFFFFFFFF) ->
with the tx sequence of 0xFFFFFFFFFFFFFFFF (18446744073709551615) with the signature ->
f87bcc405a339bf402f072c2463a19596b0a16f9b9a6ba692144398a772fd9d04164a790ea86110934f32760434afe9de1f2ac99ec52be0296629ef54394974701

0x01 = Client Identifier
0x01 = Classic Payment Instruction
0xFFFFFFFFFFFFFFFF = Amount to be sent
0xFFFFFFFFFFFFFFFF = Tx Sequence
0xC0F2329f646fd2F310862d7dB46451712201EE2e = Sender's cryptographic address
0x605302D6a2444b0dC08bf107e5FeD12f66bd0059 = Receiver's cryptographic address
f87bcc405a339bf402f072c2463a19596b0a16f9b9a6ba692144398a772fd9d04164a790ea86110934f32760434afe9de1f2ac99ec52be0296629ef54394974701 = Sender's signature of send amount, tx sequence and receiver address

So a byte location breakdown:

Byte 1 = Client Identifier
Byte 2 = Transaction Type
Bytes 2-10 = Amount to be sent
Bytes 10-18 = Tx Sequence
Bytes 18-38 = Sender Address
Bytes 38-58 = Receiver Address
Bytes 58-123 = Transaction Signature

Then it is encoded in XDR and sent through TCP