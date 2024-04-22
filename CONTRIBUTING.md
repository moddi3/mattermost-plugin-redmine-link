# Contributing to Mattermost Plugin Redmine Link

Thank you for your interest in contributing to the Mattermost Plugin Redmine Link! We welcome contributions from the community and appreciate any help in making the plugin better.

## Getting Started

1. **Fork the repository**: Fork the repository on GitHub to create a local copy of the repository.
2. **Clone the repository**: Clone the forked repository to your local machine using `git clone https://github.com/your-username/mattermost-plugin-redmine-link.git`.
3. **Build the plugin**: Run `make build` to build the plugin.

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options. In order for the below options to work, you must first enable plugin uploads via your config.json or API and restart Mattermost.

```json
    "PluginSettings" : {
        ...
        "EnableUploads" : true
    }
```

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. Edit your server configuration as follows:

```json
{
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    },
}
```

and then deploy your plugin:
```
make deploy
```

## Contributing Guidelines

1. **Create a new branch**: Create a new branch for your feature or bug fix using `git checkout -b feature/new-feature` or `git checkout -b fix/bug-fix`.
2. **Make changes**: Make the necessary changes to the code and commit them using `git add <file> && git commit -m "Added new feature"`.
3. **Test the changes**: Test the changes to ensure they do not break the plugin.
4. **Create a pull request**: Create a pull request to merge your changes into the main branch.

## Code Style

1. You can use `make check-style` to check the code style against our standards. This command runs linting checks and ensures that the code follows our formatting rules.
2. **Use consistent naming conventions**: Use consistent naming conventions for variables, functions, and classes.
3. **Write readable code**: Write code that is easy to read and understand.

## Testing

1. **Write unit tests**: Write unit tests for your code using gotestsum and jest.
2. **Run tests**: Run the tests using `make test`.

## Reporting Issues

1. **Create an issue**: Create an issue on the GitHub repository to report any bugs or issues you encounter.
2. **Provide details**: Provide as much detail as possible about the issue, including any error messages or steps to reproduce the issue.

## Code of Conduct

By participating in this project, you agree to abide by the [CODE OF CONDUCT](CODE_OF_CONDUCT.md). Please be respectful and considerate of others when interacting with the community.

## License

This plugin is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

If you have any questions or need help with contributing, please contact us at [oss@moddi3.com](mailto:oss@moddi3.com).

Thank you for your contributions!
