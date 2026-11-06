# SpectrumDF
SpectrumDF is an implementation of the [Spectrum](https://github.com/cooldogedev/Spectrum) proxy for Dragonfly servers.

## Using it
It builds against [oomph-ac/dragonfly](https://github.com/oomph-ac/dragonfly) and [oomph-ac/spectrum-internal](https://github.com/oomph-ac/spectrum-internal). Go applies `replace` lines only in the module being built, so a server that uses SpectrumDF needs the same ones in its own `go.mod`:

```
replace github.com/df-mc/dragonfly => github.com/oomph-ac/dragonfly v0.0.0-20260904153233-9221758c1eee

replace github.com/cooldogedev/spectrum => <path to spectrum-internal>
```
