# flag

is an opinionated replacement for the flag package. It aims to be a drop in replacement for 80% of the use cases
while providing an easy way for falling back to environment variables if present

## Exclude flags from env

This adds a bunch of methods *WithoutEnv which does not query environment variables.   