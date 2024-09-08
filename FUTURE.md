# Future plans for GoAddons

## Fix the updater
1. *SOLVED* ~For some reason the downloaded files are owned by ROOT:ROOT. This is not OK, they should be owned by current OS user (perhaps configurable?)~
2. *SOLVED* ~Make it non-parallel - there is no need for concurrency, there is risk for website throttling/blocking us plus extracting is fast enough as is~
3. *SOLVED* ~Make the extracting more simple - anything in the download volume that ends with .zip should be extracted and ultimately removed~
4. Look into possibility of not requiring the Docker Chrome browser when downloading files - the fewer dependencies the better

## Refactoring and generic clean-up

1. Part of the clean-up and refactoring involves more restrict code styling and that is something that is coming!
   *. Discuss the possibility of using formatters e.g. gofmt or perhaps Go linters?
2. Refactoring and making code easier to read and follow (the less jumping around the better)
3. Generic Clean-up and whatever that would entail!

## Config

1. Allow more configurations - perhaps look into some config lib that can help
2. *SOLVED* ~Look into possibly leaving JSON for the config and instead use YAML~

## Cli

1. Leave behind the current and basic CLI and look into uses something like ncurses or alternatives for it in Golang
