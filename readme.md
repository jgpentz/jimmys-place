#  Jimmys Place 🪐🧗📚

## Bumping version

1) Generate the bumpversion executable

```bash
cd bumpversion
go build
```

2) Bump the version. Make sure that all of your changes for this version are commited before this step.

> NOTE: This must be run the root project directory (until i fix it to be agnostic 👹)

```bash
# Only use one of  [major, minor, patch]
./bumpversion/bumpversion major|minor|patch
```

3) Push the changes.  The bumpversion script automatically commits the `version.json` file and creates a tag

```bash
git push && git push --tags
```
