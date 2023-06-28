# flag

is an opinionated replacement for the flag package. It aims to be a drop in replacement for 80% of the use cases
while providing an easy way for falling back to environment variables if present

## Exclude flags from env

This adds a bunch of methods *WithoutEnv which does not query environment variables.   

## Verbs

This flag package provides a list of verbs. That is something passed to the command line without a previous flag. Use `GetVerbs()` to retrieve.

    cmd verb -bool verb2 -comment text

`GetVerbs()` will return `[]string{"verb", "verb2"}`

The boolFlag will be set to true and the commentFlag will be set to "text".

## Separated

If you want to pass arguments for something like a sub command you can use the separate feature. Activate it using `flag.SetSeparate()`. Everything after `--` will be treated as a separate from command line and not parsed as verbs or flags. 

    cmd verb -bool verb2 -comment text -- separated from command line

`GetVerbs()` will return `[]string{"verb", "verb2"}`

`GetSeparate()` will return `[]string{"separated", "from", "command", "line"}`

`GetBool("bool")` will return `true`

The boolFlag will be set to true and the commentFlag will be set to "text".