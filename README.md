# shx
📜 Bash/PDKSH shell action executor. Never `npx` again.

Documentation available at [kb.ltgc.cc](https://kb.ltgc.cc/shx/).

## Supported shells
* AT&T `ksh`
  * May not work properly due to buggy `IFS` implementation.
* Bash
* `zsh`

## Unsupported shells
* `ash`
* BusyBox `sh`
* `dash`
* `loksh` (as of version 7.4)
* `oksh` (as of version 7.3)
* OpenBSD `ksh`
* `yash` (as of version 2.55)
* `mksh` (when array slicing is supported)