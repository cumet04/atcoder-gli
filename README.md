AtCoder-GLI
----------
Go implementation of [AtCoder](https://atcoder.jp/) command line tools.
This is inspired by [atcoder-cli](https://github.com/Tatamo/atcoder-cli).

### Usage

#### firststep
```
Launch a wizard for acg's initial setup for first user.
Through the wizard, you can login to atcoder in acg and setup config with descriptions.

Same function is available 'acg login' and 'acg config wizard'.

Usage:
  acg firststep [flags]

Flags:
  -h, --help   help for firststep
```

#### login
```
Login to AtCoder with USERNAME and PASSWORD.
USERNAME and PASSWORD are optional, and they are prompted if omitted.
Some actions (ex. 'acg submit') require login beforehand, so you need to login with this command.

See also 'acg help session' for current login status.

Usage:
  acg login [USERNAME] [PASSWORD] [flags]

Flags:
  -h, --help   help for login
```

#### new
```
Create new directory for CONTEST_ID and setup directories/files.
Fetch contest info from AtCoder website and download sample test cases for tasks.

Usage:
  acg new CONTEST_ID [flags]

Examples:
  For instance, created directory tree is:
  abc100/
  - .contest.json
  + a/
    - main.go // if template is set in config
    + tests/
      - sample_1.in
      - sample_1.out
      - sample_2.in
      - sample_2.out
  + b/ ...
  + c/ ...
  ...

Flags:
      --command string      command template for local test in 'acg test' (default "./{{.Script}}")
  -h, --help                help for new
      --language string     language ID specified with code submission
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'
```

#### config
```
Show/Write config values from/to config file.
Run with some config options, it write the value to file.
If you run this without any options and config file, new config file is created with default values.

See 'Global Flags' for available config options.

Usage:
  acg config [flags]
  acg config [command]

Available Commands:
  doc         Show config description with default values
  wizard      Making config wizard

Flags:
  -h, --help   help for config

Use "acg config [command] --help" for more information about a command.
```

#### config doc
```
Show config description with default values

Usage:
  acg config doc [flags]

Flags:
  -h, --help   help for doc
```

#### config wizard
```
Launch wizard for making config file, and interactively setup config parameters.

Usage:
  acg config wizard [flags]

Flags:
  -h, --help   help for wizard
```

#### lang
```
List atcoder's available languages for submit.
You can also filter languages with keyword (see 'filter' flag).

Usage:
  acg lang [flags]

Flags:
  -f, --filter string   filter keyword for list (case-insensitive)
  -h, --help            help for lang
      --no-header       Don't print header
```

#### submit
```
Submit a file as answer for a task, and wait the judge is complete.
Target file is determined by looking for a file named config's template name, in current directory.
Target task is guessed from current directory.
Language is read from config value: 'language'.

Usage:
  acg submit [flags]

Aliases:
  submit, s

Examples:
  ex 1. run in abc100/b, template = main.rb
  -> submit abc100/b/main.rb for abc100's b task

Flags:
  -h, --help     help for submit
      --nowait   exit without waiting for judge complete
```

#### session
```
Check whether current login session is alive or not.
If session is alive, it show login user's username.

See also 'acg help login'.

Usage:
  acg session [flags]

Flags:
  -h, --help   help for session
```

#### open
```
Open a contest page with default browser.
Target contest is specified by CONTEST_ID, or guessed by current directory.

See also 'acg help show' for guessing target contest specification.

Usage:
  acg open [CONTEST_ID] [flags]

Flags:
  -h, --help   help for open
```

#### show
```
Show a contest summary.
Target contest is specified by CONTEST_ID, or guessed by current directory.

If you run this command in contest directory (created by 'acg new') or under it,
target contest is guessed to the directory's contest.

If CONTEST_ID is present, use it for determining target contest (current directory is not considered).

Usage:
  acg show [CONTEST_ID] [flags]

Flags:
  -h, --help   help for show
```
