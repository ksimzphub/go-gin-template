The general template of the gin project uses zap as the log management and viper as the configuration management. By default, the common configuration items of mysql, redis and other projects are configured, and the config.yaml in the root directory can be modified.

Adding the configuration file path after using the run command makes it easier to modify the configuration file in different directories to ensure that the project can accurately read the relevant configuration.
