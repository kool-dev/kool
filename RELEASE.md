# Release process

For generating a new release we follow the steps:

- Get latest `master` branch.
- Make sure docs are updated - `kool run make-docs`
- Make sure presets are updated - `kool run parse-presets`
- Make sure formatting is correct - `kool run fmt`
- Make sure tests are passing - `kool run test`
- Pick the version name you wanna build - `export KOOL_VERSION=0.0.0` (taking into consideration [Semantic Versioning rules for Major, Minor and Patch versions](https://semver.org/#summary))
- Build all artifacts - `bash build_artifacts.sh`
- Draft new release
- Go through all PRs merged and write the changelog (Added / Fixed / Removed / Changed / Notes)
- Upload to the release all the artifacts built at dist/ folder.
- Publish the new version (which will create the tag as well)
- In case of version bumping, check if we need to update `SECURITY.md` to show what version we currently support.
