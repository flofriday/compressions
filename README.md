# compressions
Playing with compression algorithms

## LZSS
Im progress, but compressing does work at the moment somewhat

This implementation uses 12 bit for the lookback and 4 bit for the length.
At the moment it uses sequential search.
Beginning with the 1st byte, every eigth byte is a bit-mask to indicate which 
bytes are values and which bytes are the pointers to previous values.

### Run the code
```
go build lzss.go
./lzss
```

### Resources
http://wiki.xentax.com/index.php/LZSS
http://michael.dipperstein.com/lzss/ 