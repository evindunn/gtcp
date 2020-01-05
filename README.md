# GoLang TCP

<a href="https://github.com/evindunn/Tcp/actions?query=workflow%3ABuild">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/Tcp/workflows/Build/badge.svg">
</a>

<a href="https://github.com/evindunn/Tcp/actions?query=workflow%3ATest">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/Tcp/workflows/Test/badge.svg">
</a>

<a href='https://coveralls.io/github/evindunn/Tcp?branch=master'>
  <img src='https://coveralls.io/repos/github/evindunn/Tcp/badge.svg?branch=master&service=github' alt='Coverage Status' />
</a>



A message-passer built on TCP written in GO
- zlib compression for larger messages
- Utilities for converting the Tcp.Message to and from raw bytes
- Utilities for reading/writing Tcp.Message to/from net.Conn

```
|------- 8 bytes---------|------- 1 byte ---------|------- Remaining bytes ---------|
|----- messageSize ------|----- isCompressed -----|---------- content --------------|
```