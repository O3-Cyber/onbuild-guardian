# Onbuild Guardian

ONBUILD directives pose a supply chain risk. Onbuild Guardian lets you scan for them, and allow-list the ones you're certain are safe

Installation: `go install github.com/O3-Cyber/onbuild-guardian@latest`
Usage: `onbuild-guardian <base-image>`

When you identify ONBUILD directives that you're ok with, add them to `.onbuild-guardian.yaml` like so:

```yaml
# allow the specific ONBUILD directives. Onbuild Guardian will exit with a non-zero status code if there's an ONBUILD directives that isn't on the allow list
allowed:
  # first, the image that we're allowing the directives in
  node:onbuild:
    # then the directives themselves
    - ENV ...
    - ARG ...
    - RUN ...

# allow all ONBUILD directives for a given image
allowAll:
  python:onbuild: true
```
