# Hydrogen Lib
### THIS LIBRARY IS FOR NODE TO NODE COMMUNICATION ONLY
### WHEN USING THE NETWORK AS A USER OF THE NETWORK
### USE NORMAL (g)RPC TO INTERACT WITH THE NETWORK
This is a data transport lib that works on barebones TCP (transmitting bytes only), but has cross-compatability
with HTTP(s) and RPC

### Types supported by Hydrogen
The supported types goes as follows:

- string (Identifier "str")
- i32 and u32
- i64 and u64
- i128 and u128

(You can provide complex structs and more data types if you provide an xdr encoded string)
