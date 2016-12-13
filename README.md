# Popcorn

Go standard library makes file manipulation very easy. However, it's a little
cumbersome to generate a whole file architecture easily. This `pop` Go package
is here to help! Simply describe your tree with a few *pop.Corn*, call
`pop.Generate` or `pop.GenerateFromRoot` on it and all the intermediate
directories and dummy files will be ready in no time! Populating a directory
for integration testing has never been so fast and fun.
