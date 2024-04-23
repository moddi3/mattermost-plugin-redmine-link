# Redmine Link Transform for Mattermost · ![ci workflow](https://github.com/moddi3/mattermost-plugin-redmine-link/actions/workflows/ci.yml/badge.svg) [![Mutable.ai Auto Wiki](https://img.shields.io/badge/Auto_Wiki-Mutable.ai-blue)](https://wiki.mutable.ai/moddi3/mattermost-plugin-redmine-link)

This plugin enhances the Mattermost experience by transforming Redmine links in messages into rich content. It automatically extracts information from Redmine issues and presents them in an easy-to-read format.

Automatically transforms Redmine issue links in Mattermost messages into a readable format, providing additional information such as issue status, priority, assignee, and other relevant details directly in the chat.

## Usage

Include Redmine issue links in your Mattermost messages to see the plugin in action. The plugin needs to be configured before use.

### Example

Here's an example of how the plugin transforms Redmine links:

- **Original message**: "Check out the issue here: https://www.redmine.org/issues/3451"
- **Transformed message**: "Check out the issue here: [Defect#3451: Issue Creation Via Email not Working ](https://www.redmine.org/issues/3451 "Assignee: Unassigned&#013;Priority: Normal&#013;Status: Closed&#013;Author: Carlo Camerino&NewLine;Last update: Fri, 16 Jul 2021 15:42:40 EEST")"

In the transformed message, you can view the issue subject and tracker. Additional details such as status, priority, author, and last update date are available when hovering over the link.
#### Example of a transfrmed message in Markdown format:
```
[Defect#3451: Issue Creation Via Email not Working ](https://www.redmine.org/issues/3451 "Assignee: Unassigned&#013;Priority: Normal&#013;Status: Closed&#013;Author: Carlo Camerino&#013;Last update: Fri, 16 Jul 2021 15:42:40 EEST")
```
Carriage Return `&#013;` is used to insert a newline character in the link's title attribute. This allows the additional information (Assignee, Priority, Status, Author, Last update) to be displayed on separate lines when the user hovers over the link or views it in a Markdown renderer that supports tooltips.

## Installation

While you have the option to build the plugin yourself, it is much easier to download the already built plugin from the [releases page](https://github.com/moddi3/mattermost-plugin-redmine-link/releases) of the GitHub repository. Once downloaded, follow the Mattermost documentation on [plugin installation](https://developers.mattermost.com/integrate/plugins/components/server/hello-world/#install-the-plugin) to install the plugin in your Mattermost server.

### Building Package

1. Clone the repository: `git clone https://github.com/moddi3/mattermost-plugin-redmine-link`
2. Build the plugin: `make build`
3. Follow the Mattermost documentation on [plugin installation](https://developers.mattermost.com/integrate/plugins/components/server/hello-world/#install-the-plugin) to install the plugin in your Mattermost server.

## Configuration

After installation, configure the plugin in the Mattermost System Console:

- **Redmine Instance URL**: Specify the URL of your Redmine instance.
- **Redmine API Key (optional)**: Add your Redmine API key to allow the plugin to fetch issue data (only if you are using private redmine instance).

## Documentation

For more detailed documentation and usage instructions, visit the [wiki page](https://wiki.mutable.ai/moddi3/mattermost-plugin-redmine-link).

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING](CONTRIBUTING.md) file for guidelines.

## License

This plugin is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Mattermost Starter Plugin Info
For information on getting started with Mattermost plugins, refer to the [PLUGIN.md](PLUGIN.md) file in the root directory of this repository.

## Support

If you encounter any issues or have any questions, please raise them in the [GitHub repository](https://github.com/moddi3/mattermost-plugin-redmine-link) issues section.
