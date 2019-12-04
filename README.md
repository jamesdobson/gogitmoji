# gogitmoji

[Gitmoji](https://gitmoji.carloscuesta.me/) helper written with ❤️ in Go. Inspired by [gitmoji-cli](https://github.com/carloscuesta/gitmoji-cli).

Currently under construction 🚧, not much is implemented!

[![Build Status](https://travis-ci.org/jamesdobson/gogitmoji.svg?branch=master)](https://travis-ci.org/jamesdobson/gogitmoji)
[![Go Report Card](https://goreportcard.com/badge/github.com/jamesdobson/gogitmoji)](https://goreportcard.com/report/github.com/jamesdobson/gogitmoji)
[![Coverage Status](https://coveralls.io/repos/github/jamesdobson/gogitmoji/badge.svg?branch=master)](https://coveralls.io/github/jamesdobson/gogitmoji?branch=master)

[![asciicast](https://asciinema.org/a/283799.svg)](https://asciinema.org/a/283799)

## Usage

```bash
gitmoji help
```

```
gogitmoji helps you write git commit messages containing gitmoji!

Usage:
  gogitmoji [flags]
  gogitmoji [command]

Available Commands:
  commit      ⚡️  Compose a commit message and execute git commit
  help        📗  Help about any command
  info        🌍  Open gimoji information page in your browser
  list        📜  List all available gitmoji (default command)
  update      🔄  Update the list of gitmoji

Flags:
      --config string   config file (default is $HOME/.gogitmoji.yaml)
  -h, --help            help for gogitmoji

Use "gogitmoji [command] --help" for more information about a command.
```

### Commit

Guides the user through the process of composing a commit message, and then
executes `git commit`.

```bash
gitmoji commit
```

`commit` is the default command, so just the following is equivalent:

```bash
gitmoji
```

### List

Prints the list of gitmoji.

```bash
gitmoji list
```

```
🎨  - :art: Improving structure / format of the code.
⚡️  - :zap: Improving performance.
🔥  - :fire: Removing code or files.
🐛  - :bug: Fixing a bug.
🚑  - :ambulance: Critical hotfix.
✨  - :sparkles: Introducing new features.
📝  - :pencil: Writing docs.
🚀  - :rocket: Deploying stuff.
💄  - :lipstick: Updating the UI and style files.
🎉  - :tada: Initial commit.
✅  - :white_check_mark: Updating tests.
🔒  - :lock: Fixing security issues.
🍎  - :apple: Fixing something on macOS.
🐧  - :penguin: Fixing something on Linux.
🏁  - :checkered_flag: Fixing something on Windows.
🤖  - :robot: Fixing something on Android.
🍏  - :green_apple: Fixing something on iOS.
🔖  - :bookmark: Releasing / Version tags.
🚨  - :rotating_light: Removing linter warnings.
🚧  - :construction: Work in progress.
💚  - :green_heart: Fixing CI Build.
⬇️  - :arrow_down: Downgrading dependencies.
⬆️  - :arrow_up: Upgrading dependencies.
📌  - :pushpin: Pinning dependencies to specific versions.
👷  - :construction_worker: Adding CI build system.
📈  - :chart_with_upwards_trend: Adding analytics or tracking code.
♻️  - :recycle: Refactoring code.
🐳  - :whale: Work about Docker.
➕  - :heavy_plus_sign: Adding a dependency.
➖  - :heavy_minus_sign: Removing a dependency.
🔧  - :wrench: Changing configuration files.
🌐  - :globe_with_meridians: Internationalization and localization.
✏️  - :pencil2: Fixing typos.
💩  - :poop: Writing bad code that needs to be improved.
⏪  - :rewind: Reverting changes.
🔀  - :twisted_rightwards_arrows: Merging branches.
📦  - :package: Updating compiled files or packages.
👽  - :alien: Updating code due to external API changes.
🚚  - :truck: Moving or renaming files.
📄  - :page_facing_up: Adding or updating license.
💥  - :boom: Introducing breaking changes.
🍱  - :bento: Adding or updating assets.
👌  - :ok_hand: Updating code due to code review changes.
♿️  - :wheelchair: Improving accessibility.
💡  - :bulb: Documenting source code.
🍻  - :beers: Writing code drunkenly.
💬  - :speech_balloon: Updating text and literals.
🗃  - :card_file_box: Performing database related changes.
🔊  - :loud_sound: Adding logs.
🔇  - :mute: Removing logs.
👥  - :busts_in_silhouette: Adding contributor(s).
🚸  - :children_crossing: Improving user experience / usability.
🏗  - :building_construction: Making architectural changes.
📱  - :iphone: Working on responsive design.
🤡  - :clown_face: Mocking things.
🥚  - :egg: Adding an easter egg.
🙈  - :see_no_evil: Adding or updating a .gitignore file
📸  - :camera_flash: Adding or updating snapshots
⚗  - :alembic: Experimenting new things
🔍  - :mag: Improving SEO
☸️  - :wheel_of_dharma: Work about Kubernetes
🏷️  - :label: Adding or updating types (Flow, TypeScript)
🌱  - :seedling: Adding or updating seed files
🚩  - :triangular_flag_on_post: Adding, updating, or removing feature flags
💫  - :dizzy: Adding or updating animations and transitions
```

### Update

Checks to see if there is a new list of gitmoji online, updating the local cache
if there are new gitmoji.

```bash
gitmoji update
```

## Capabilities

- [x] Select a gitmoji from a list
- [x] Search name, description, and code
- [x] Automatically get latest gitmoji database
- [x] Implement first automated tests
- [x] Ask gitmoji commit message questions and output commit message
- [x] Commit to git with commit message
- [x] List gitmoji
- [x] Update gitmoji cache
- [x] Support emoji code in message instead of emoji itself (add configuration support)
- [x] Ask for scope depending on config
- [x] CI build and releases
- [ ] Homebrew package
- [ ] Install git commit hook
- [ ] Support conventional commits (with optional gitmoji)

## License

Licensed under the [MIT](https://github.com/jamesdobson/gogitmoji/blob/master/LICENSE) license.

The gitmoji database is from [Gitmoji](https://gitmoji.carloscuesta.me/).
