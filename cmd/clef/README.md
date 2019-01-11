Clef
----
Clef can be used to sign transactions and data and is meant as a replacement for gvnt's account management.
This allows DApps not to depend on gvnt's account management. When a DApp wants to sign data it can send the data to
the signer, the signer will then provide the user with context and asks the user for permission to sign the data. If
the users grants the signing request the signer will send the signature back to the DApp.
  
This setup allows a DApp to connect to a remote VNT node and send transactions that are locally signed. This can
help in situations when a DApp is connected to a remote node because a local VNT node is not available, not
synchronised with the chain or a particular VNT node that has no built-in (or limited) account management.
  
