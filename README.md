# RSA-Shamir-Secret-Sharing

[![GoDoc](https://godoc.org/github.com/ceriath/rsa-shamir-secret-sharing?status.svg)](https://godoc.org/github.com/ceriath/rsa-shamir-secret-sharing)
[![Go Report Card](https://goreportcard.com/badge/github.com/ceriath/rsa-shamir-secret-sharing)](https://goreportcard.com/report/github.com/ceriath/rsa-shamir-secret-sharing)
[![Build Status](https://travis-ci.com/ceriath/rsa-shamir-secret-sharing.svg?branch=master)](https://travis-ci.com/ceriath/rsa-shamir-secret-sharing)

---

## **This tool has been written for academic reasons only. I do not recommend to use it in production since it may contain security flaws. It is written with the best of my knowledge and belief, but i am not a cryptographer and can not guarantee anything.**  
However, if you find any security flaws, please open an issue without disclosing the flaw or contact me via email. I prefer responsible disclosures.  

---

RSA-Shamir-Secret-Sharing is a little tool i wrote for my Bachelor-Thesis. It implements the [Shamir-Secret-Sharing-Scheme](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing). 
It is especially focused on sharing RSA key pairs, but also works perfectly fine with everything else.

---

## Usage

### rsasss
```
> rsasss -h

Implements the shamir secret sharing especially made for rsa private keys

Usage:
  rsasss [flags]
  rsasss [command]

Available Commands:
  gen         Generates a shared RSA KeyPair
  help        Help about any command
  join        Joins shares to reconstruct the secret
  split       Splits a given secret into shares
  version     Print the version number

Flags:
  -h, --help   help for rsasss

Use "rsasss [command] --help" for more information about a command.
```

---

### rsasss gen

```
> rsasss gen -h
Generates a public RSA key and n shares for the private key of which k are required for reconstruction

Usage:
  rsasss gen [flags]

Flags:
  -b, --bitsize int      bitsize of the generated key (default 2048)
  -c, --createKeyfile    set to true if you want to save the private key to a file
  -h, --help             help for gen
  -f, --saveShares       set to true if you want to save the shares to a file
  -n, --sharecount int   number of shares to generate (default 5)
  -k, --threshold int    number of required shares to reconstruct the secret (default 3)
```

---

### rsasss join

```
> rsasss join -h
Joins multiple shares to reconstruct a secret

Usage:
  rsasss join [flags]

Flags:
  -h, --help   help for join

```

---

### rsasss split

```
> rsasss split -h
Splits the given secret into n shares of which k are required to reconstruct the secret

Usage:
  rsasss split [flags]

Flags:
  -h, --help                help for split
  -f, --saveShares          set to true if you want to save the shares to a file
  -s, --secret string       the secret to split
  -p, --secretpath string   filepath to the secret
  -n, --sharecount int      number of shares to generate (default 5)
  -k, --threshold int       number of required shares to reconstruct the secret (default 3)
```

---

## Contribution
I am not really developing this library and tool actively.  
If you have any issues with it, feel free to open an issue or a pull request and i'll look into it. 

---

## Structure
In [cmd](./cmd) you can find the sources that are required to run the CLI tool  
In [gf256](./gf256) you can find an implementation of finite fields (written by [@rsc](https://github.com/rsc) and provided under the BSD-3-License)  
in [s4](./s4) you can find my implementation of the **S**hamir-**S**ecret-**S**haring-**S**cheme (s4)  

---

## Legal

This work is licensed under the MIT License (see [LICENSE.md](./LICENSE.md))  
The GF256 Package is licensed under the BSD-3-License and belongs to [@rsc](https://github.com/rsc)