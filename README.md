# gogitmoji

[Gitmoji](https://gitmoji.carloscuesta.me/) helper written with ❤️ in Go. Inspired by [gitmoji-cli](https://github.com/carloscuesta/gitmoji-cli).

[![Build Status](https://travis-ci.org/jamesdobson/gogitmoji.svg?branch=master)](https://travis-ci.org/jamesdobson/gogitmoji)
[![Go Report Card](https://goreportcard.com/badge/github.com/jamesdobson/gogitmoji)](https://goreportcard.com/report/github.com/jamesdobson/gogitmoji)
[![Coverage Status](https://coveralls.io/repos/github/jamesdobson/gogitmoji/badge.svg?branch=master)](https://coveralls.io/github/jamesdobson/gogitmoji?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/386deea133c0488a88a04b3bb1c44244)](https://www.codacy.com/manual/jamesdobson/gogitmoji?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jamesdobson/gogitmoji&amp;utm_campaign=Badge_Grade)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjamesdobson%2Fgogitmoji.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjamesdobson%2Fgogitmoji?ref=badge_shield)

[![asciicast](https://asciinema.org/a/321491.svg)](https://asciinema.org/a/321491)

## Installation

### Mac OS X

Install using Homebrew:

```bash
brew install jamesdobson/gogitmoji/gogitmoji
```

## Usage

```bash
gitmoji help
```

```console
gogitmoji helps you write git commit messages containing gitmoji!

Usage:
  gitmoji [flags]
  gitmoji [command]

Available Commands:
  commit      ⚡️  Compose a commit message and execute git commit (default command)
  export      🚢  Export a commit template
  help        📗  Help about any command
  hook        🎣  Manage commit hooks
  info        🌍  Open gimoji information page in gyour browser
  list        📜  List all available gitmoji
  update      🔄  Update the list of gitmoji
  version     ℹ️  Display the version of this program

Flags:
      --config string   config file (default is $HOME/.gitmoji/config.yaml)
  -h, --help            help for gitmoji

Use "gitmoji [command] --help" for more information about a command.
```

### Commit

Guides the user through the process of composing a commit message, and then
executes `git commit`.

```console
gitmoji commit
```

`commit` is the default command, so just the following is equivalent:

```console
gitmoji
```

### Git Hook

You can configure git to run gogitmoji automatically when you execute `git commit`,
so that you don't always have to remember to type `gitmoji`. To do so, just
set up a git hook in your git repositories:

```bash
#!/bin/sh

exec < /dev/tty
gitmoji hook do $1
exit $?
```

Place the above in the file `<my repo>/.git/hooks/prepare-commit-msg` and ensure
the file is executable (e.g. `chmod a+x prepare-commit-msg`)

### List

Prints the list of gitmoji.

```console
gitmoji list
```

```console
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
🥅  - :goal_net: Catching errors
💫  - :dizzy: Adding or updating animations and transitions
🗑  - :wastebasket: Deprecating code that needs to be cleaned up.
```

### Update

Checks to see if there is a new list of gitmoji online, updating the local cache
if there are new gitmoji.

```console
gitmoji update
```

## Configuration

The configuration file is stored at `~/.gitmoji/config.yaml`. The config file
can specify the following:

- Default commit template
- Enable "scope" prompt
- Emoji format
- New commit templates

### Set the Default Commit Template

Specify the name of the default commit template:

```yaml
template: "conventional"
```

### Enable the "scope" Prompt

Enable / disable prompting for commit scope:

```yaml
scope: True
```

Note: this is used by the default `gitmoji` template, but has no effect on the
default `conventional` template. This can be changed by defining a custom
template and using the `enablesetting` field on the corresponding `Prompt`.

### Set the Emoji Format

The emoji format can be set to either `emoji` (its default value) or `code`:

```yaml
format: code
```

When set to `emoji`, the UTF-8 bytes encoding the emoji will be used in the
commit message. When set to `code`, an text string (e.g. `:sparkles:`) will be
used. GitHub will render this as an emoji.

### Define New Commit Templates

The configuration file allows the definition of new commit templates. A commit
template has three parts:

- The command to run
- Arguments to the command
- User prompts of inputs for the arguments

The following example demonstrates all three elements:

```yaml
templates:
  example:
    Command: echo
    CommandArgs:
    - Hello, {{ .name }}, I'm pleased to meet you.
    Prompts:
    - prompttype: text
      mandatory: true
      prompt: Enter your name
      valuecode: name
```

This example prompts the user to enter their name, and then echoes it back
to them as a polite greeting:

```console
jamesdobson@MacBook-Pro gitmoji % gitmoji commit -t example
Using config file: /Users/jamesdobson/.gitmoji/config.yaml
✔ Enter your name: James
Going to execute: echo "Hello, James, I'm pleased to meet you."

Execute: y
Executing...
Hello, James, I'm pleased to meet you.

gogitmoji done.
```

In this example, the command is `echo`. In most cases, however, `Command` should
be set to `git`.

The arguments to the command are expressed as an array under the `CommandArgs`
section. Each argument is a [Go template](https://golang.org/pkg/text/template/)
that can refer to inputs that come from the user prompts. If an argument evaluates
to the empty string, it is skipped.

The final section, `Prompts`, is an array of user prompts. There are 3 kinds of
user prompt, differentiated by their `Type` field:

- `text`: Prompts the user with the text in `Prompt`, and waits for the user to enter a text response.
- `choice`: Prompts the user with a selection of options as given by the `Choices` field.
- `gitmoji`: Prompts the user with a list of gitmoji.

The result of the prompt is stored under the name given by the `Name` field and
is made available in the command arguments via the `{{ .xyz }}` syntax, where
`xyz` is whatever was specified in the `Name` field.

There is an additional section, `Messages`, that is used when gogitmoji is called
as a commit hook. In this case, no command is executed (because commit is already
running) however the `Messages` are evaluated and written to the file that git
provides to the commit hook as an argument.

#### Default gitmoji commit template

This is the default `gitmoji` commit template:

```yaml
templates:
  gitmoji:
    Command: git
    CommandArgs:
    - commit
    - -m
    - '{{if eq (getString "format") "emoji"}}{{.gitmoji.Emoji}} {{else}}{{.gitmoji.Code}}{{end}}
      {{with .scope}}({{.}}): {{end}}{{.title}}'
    - '{{with .message}}-m{{end}}'
    - '{{.message}}'
    Messages:
    - '{{if eq (getString "format") "emoji"}}{{.gitmoji.Emoji}} {{else}}{{.gitmoji.Code}}{{end}}
      {{with .scope}}({{.}}): {{end}}{{.title}}'
    - '{{.message}}'
    Prompts:
    - Type: gitmoji
      Mandatory: true
      Name: gitmoji
    - Type: text
      Prompt: Enter the scope of current changes
      Name: Scope
      Condition: scope
    - Type: text
      Mandatory: true
      Prompt: Enter the commit title
      Name: title
    - Type: text
      Prompt: Enter the (optional) commit message
      Name: message
```

#### Default conventional commit template

This is the default `conventional` commit template:

```yaml
templates:
  conventional:
    Command: git
    CommandArgs:
    - commit
    - -m
    - '{{.type}}: {{.description}}'
    - '{{with .body}}-m{{end}}'
    - '{{.body}}'
    - '{{with .footer}}-m{{end}}'
    - '{{.footer}}'
    Messages:
    - '{{.type}}: {{.description}}'
    - '{{.body}}'
    - '{{.footer}}'
    Prompts:
    - Type: choice
      Mandatory: true
      Prompt: 'Choose the type of commit:'
      Name: type
      Choices:
      - Value: feat
        Description: A new feature.
      - Value: fix
        Description: A bug fix.
      - Value: docs
        Description: Documentation only changes.
      - Value: perf
        Description: A code change that improves performance.
      - Value: refactor
        Description: A code change that neither fixes a bug nor adds a feature.
      - Value: test
        Description: Adding missing or correcting existing tests.
      - Value: chore
        Description: Changes to the build process or auxiliary tools and libraries
          such as documentation generation.
    - Type: text
      Mandatory: true
      Prompt: Enter the commit description, with JIRA number at end
      Name: description
    - Type: text
      Prompt: Enter the (optional) commit body
      Name: body
    - Type: text
      Prompt: Enter the (optional) commit footer
      Name: footer
```

## License

Licensed under the [MIT](https://github.com/jamesdobson/gogitmoji/blob/master/LICENSE) license.

The gitmoji database is from [Gitmoji](https://gitmoji.carloscuesta.me/).

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjamesdobson%2Fgogitmoji.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjamesdobson%2Fgogitmoji?ref=badge_large)
