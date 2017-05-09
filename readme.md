NTGO
====
[![Build Status](https://travis-ci.org/HowardStark/ntgo.svg?branch=master)](https://travis-ci.org/HowardStark/ntgo)

An implementation of the FRC key-value store networking protocol "NetworkTables". NTGo is aiming for full NetworkTables 3.0 compliance. As of right now, there is no plan for 2.0 backwards compatibility (the protocol is itself backwards compatible but NTGo will currently reject any 2.0 client).

## Roadmap

### Done
- Support for all entry types (including arrays)
- Decoding of all messages (excluding RPC)
- Basic Server/Client Architecture

### Todo
- Message handling
- Initialization message flow
- Persistent caching
    - Caching abstraction to allow for custom caching mechanisms without code change
- RPC Support

## Questions

If you have any questions about the project, feel free to email me at howard@getcoffee.io





