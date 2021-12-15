# Hide
A simple steganographic tool to hide byte streams in pictures.

## Example executions
### Encoding a byte stream
```bash
.\hide.exe -image=".\resources\input.png" -payload=".\resources\swag-cat.png" -output=".\encoded"      
```

### Decoding a byte stream from the generated image
```bash
.\hide.exe -mode="decode" -image=".\encoded.png" -output=".\decoded.png"
```