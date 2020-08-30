AtCoder-GLI
----------
Go implementation of [AtCoder](https://atcoder.jp/) command line tools.
This is inspired by [atcoder-cli](https://github.com/Tatamo/atcoder-cli).

### TODO
* config document
### Usage

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

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'
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
    + samples/
      - sample_1.in
      - sample_1.out
      - sample_2.in
      - sample_2.out
  + b/ ...
  + c/ ...
  ...

Flags:
  -h, --help   help for new

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
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
  lang        Select preferred language for submit

Flags:
  -h, --help   help for config

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'

Use "acg config [command] --help" for more information about a command.
```

#### config lang
```
Search and select preferred language.
Selected language is saved in config and used when submit code.

You can search language with keyword (prompted) and choose one from them.
Search targets are all available languages in AtCoder,
and keyword is case-insensitive.

Usage:
  acg config lang [flags]

Flags:
  -h, --help   help for lang

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'
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

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'
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

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'
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

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'
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

Global Flags:
      --command string      command template that runs in 'acg test' (default "./{{.ScriptFile}}")
      --language string     language id used as submit code's language
      --sample_dir string   directory name where sample in/out files are stored in (default "tests")
      --template string     template file name that is copied to task directory in 'acg new'
```
