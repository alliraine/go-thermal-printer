# ESC/POS image pipeline performance

The `ImageToRasterBytes` conversion was benchmarked before and after replacing the
manual scaling and grayscale loops with the optimized `golang.org/x/image/draw`
implementation.

## Benchmark setup

```
go test ./pkg/escpos -bench=ImageToRasterBytes -benchmem
```

Benchmark images are procedurally generated gradients sized to represent
common receipt use cases:

- Small: 384×256px, already within printer width
- Medium: 1024×768px, typical photo snapped on a phone
- Large: 2048×1536px, a down-scaled camera image

## Results

| Implementation | Small (ns/op) | Medium (ns/op) | Large (ns/op) | Allocations (B/op) |
| -------------- | ------------- | -------------- | ------------- | ------------------ |
| Previous manual loops | 5,375,174 | 5,987,783 | 5,650,285 | ~1.0 MB / 98k allocs |【e3e7bc†L1】【7fa902†L1】【f8f5ce†L1-L3】
| `draw.ApproxBiLinear` + `image/draw` | 1,606,471 | 6,165,524 | 6,189,398 | ~0.59 MB / 7 allocs |【4185d8†L1-L5】【5f4aca†L1】【f533b7†L1-L3】

## Takeaways

- Scaling now uses `draw.ApproxBiLinear`, producing smoother results and cutting
  allocations by two orders of magnitude thanks to the more efficient pipeline.
- Runtime for already-on-width assets dropped by ~70%, and larger images retain
  comparable throughput while delivering higher quality resampling.
