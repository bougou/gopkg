# go x509


```
                pem.Decode
+-----------+  +-----------> +-----------+
| PEM Bytes |                | pem.Block |
+-----------+  <-----------+ +-----------+
                pem.Encode
                                   +
                                   v .Bytes
                   ParseCertificate
+------------------+ <-----+ +------------+
| x509.Certificate |         | DER Bytes  |
+------------------+ +-----> +------------+
                       .Raw

```
