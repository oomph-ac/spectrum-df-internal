# SpectrumDF
SpectrumDF is an implementation of the [Spectrum](https://github.com/cooldogedev/Spectrum) proxy for Dragonfly servers.

## Using it
It builds against [oomph-ac/dragonfly](https://github.com/oomph-ac/dragonfly). Go applies `replace` lines only in the module being built, so a server that uses SpectrumDF needs the same one in its own `go.mod`:

```
replace github.com/df-mc/dragonfly => github.com/oomph-ac/dragonfly v0.0.0-20260904153233-9221758c1eee
```

Spectrum's frame reader and writer, and its own packets such as `Transfer`, are copied in from [oomph-ac/spectrum-internal](https://github.com/oomph-ac/spectrum-internal), so a server does not need Spectrum itself. Code that imported them from Spectrum imports them from here at the same path, with `github.com/cooldogedev/spectrum/` changed to `github.com/cooldogedev/spectrum-df/`, for example `github.com/cooldogedev/spectrum-df/server/packet`.
