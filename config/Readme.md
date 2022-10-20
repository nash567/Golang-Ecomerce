# Configuration for application

# steps

- first we will read env files and load env variables
- then we will read our config files using the package configor 

# Working of Configor package
- this package first checks all the environment variables with a given prefix then it will check the files for the variable if the variable is present in env then it will take the variables value and embed into the given Config struct if the variable is not present in environment then it will take the variable from the given config file and embed it into the Config struct

- it does not override the environment variables values over the config file value